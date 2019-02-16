package job

import (
	"context"
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
		log.Panic("Unable to fetch pull requests:", err)
	}

	for _, pr := range pullRequests {
		reviews, resp, err := rf.client.PullRequests.ListReviews(context.Background(), pr.R.Project.RepoOwner, pr.R.Project.RepoName, int(pr.Number), nil)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Panic("Unable to fetch pull request reviews:", err)
		}
		log.Printf("(%d) review(s) fetched for pull request %q\n", len(reviews), pr.Title)

		reviewers := make([]*models.User, 0)
		for _, review := range reviews {
			user := transformUser(review.GetUser(), rf.storage.Get())
			reviewers = append(reviewers, user)

			if review.GetState() == approvedState {
				exists, err := user.ApprovedPullRequests(qm.Where("pull_request_id = ?", pr.ID)).Exists(context.Background(), rf.storage.Get())
				if err != nil {
					log.Panic("Unable to check pull request activity record:", err)
				}
				if !exists {
					err = user.AddApprovedPullRequests(context.Background(), rf.storage.Get(), false, pr)
					if err != nil {
						log.Panic("Unable to persist user approved pull request")
					}
				}
			} else {
				exists, err := user.CommentedPullRequests(qm.Where("pull_request_id = ?", pr.ID)).Exists(context.Background(), rf.storage.Get())
				if err != nil {
					log.Panic("Unable to check pull request activity record:", err)
				}
				if !exists {
					err = user.AddCommentedPullRequests(context.Background(), rf.storage.Get(), false, pr)
					if err != nil {
						log.Panic("Unable to persist user commented pull request")
					}
				}
			}
			if err != nil {
				log.Panic("Unable to persist pull request activity record:", err)
			}
		}

		idlers := findIdlers(pr.R.Reviewers, reviewers)
		pr.AddIdlers(context.Background(), rf.storage.Get(), false, idlers...)

		pr.Update(context.Background(), rf.storage.Get(), boil.Infer()) // updates the 'updated_at' column
	}
}

func findIdlers(requestedReviewers, actualReviewers []*models.User) []*models.User {
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
