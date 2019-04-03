package job

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/google/go-github/github"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"log"
	"net/http"
)

const approvedState = "APPROVED"

// ReviewFetcher is a regular job which fetches reviews for tracked pull requests
type ReviewFetcher struct {
	storage storage.Storage
	client  *github.Client
}

// NewReviewFetcher creates an instance of ReviewFetcheer
func NewReviewFetcher(storage storage.Storage, client *github.Client) *ReviewFetcher {
	return &ReviewFetcher{
		storage: storage,
		client:  client,
	}
}

// Name returns the name of the ReviewFetcher job
func (rf *ReviewFetcher) Name() string {
	return reviewFetcher
}

// Period returns the period of time when the ReviewFetcher job should execute
func (rf *ReviewFetcher) Period() string {
	return "@every 15m"
}

// Run executes the ReviewFetcher job
func (rf *ReviewFetcher) Run() {
	log.Printf("STARTING %v job", rf.Name())
	defer log.Printf("FINISHED %v job", rf.Name())

	pullRequests, err := models.PullRequests(qm.Load("Project"), qm.Load("Reviewers")).All(context.Background(), rf.storage.Get())
	if err != nil {
		log.Println("Unable to fetch pull requests:", err)
		return
	}

	for _, pullRequest := range pullRequests {
		reviews, resp, err := rf.client.PullRequests.ListReviews(context.Background(), pullRequest.R.Project.RepoOwner, pullRequest.R.Project.RepoName, int(pullRequest.Number), nil)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Println("Unable to fetch pull request reviews:", err)
			continue
		}
		log.Printf("(%d) review(s) fetched for pull request %q\n", len(reviews), pullRequest.Title)

		txErr := rf.storage.Transaction(context.Background(), func(context context.Context, tx *sql.Tx) error {
			reviewers := make([]*models.User, 0)
			for _, review := range reviews {
				user, err := transformUser(context, tx, review.GetUser())
				if err != nil {
					return fmt.Errorf("unable to transform reviewer user: %s", err)
				}

				reviewers = append(reviewers, user)

				if review.GetState() == approvedState {
					exists, err := user.ApprovedPullRequests(qm.Where("pull_request_id = ?", pullRequest.ID)).Exists(context, tx)
					if err != nil {
						return err
					}
					if !exists {
						err := user.AddApprovedPullRequests(context, tx, false, pullRequest)
						if err != nil {
							return err
						}
					}
				} else {
					exists, err := user.CommentedPullRequests(qm.Where("pull_request_id = ?", pullRequest.ID)).Exists(context, tx)
					if err != nil {
						return err
					}
					if !exists {
						err := user.AddCommentedPullRequests(context, tx, false, pullRequest)
						if err != nil {
							return err
						}
					}
				}
			}

			idlers := rf.findIdlers(pullRequest.R.Reviewers, reviewers)
			err = pullRequest.AddIdlers(context, tx, false, idlers...)
			if err != nil {
				return err
			}

			_, err = pullRequest.Update(context, tx, boil.Infer()) // updates the 'updated_at' column
			return err
		})

		if txErr != nil {
			log.Printf("Unable to persist pull request %q activity records: %v\n", pullRequest.Title, txErr)
		}
	}
}

func (rf *ReviewFetcher) findIdlers(requestedReviewers, actualReviewers []*models.User) []*models.User {
	idlers := make([]*models.User, 0)
	for _, requestedReviewer := range requestedReviewers {
		if !contains(actualReviewers, requestedReviewer) {
			idlers = append(idlers, requestedReviewer)
		}
	}
	return idlers
}

func contains(users []*models.User, user *models.User) bool {
	for _, u := range users {
		if u.GithubID == user.GithubID {
			return true
		}
	}
	return false
}
