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

	db := rf.storage.Get()

	pullRequests, err := models.PullRequests(qm.Load("Project"), qm.Load("Reviewers")).All(context.Background(), db)
	if err != nil {
		log.Println("Unable to fetch pull requests:", err)
		return
	}

	for _, pr := range pullRequests {
		reviews, resp, err := rf.client.PullRequests.ListReviews(context.Background(), pr.R.Project.RepoOwner, pr.R.Project.RepoName, int(pr.Number), nil)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Println("Unable to fetch pull request reviews:", err)
			continue
		}
		log.Printf("(%d) review(s) fetched for pull request %q\n", len(reviews), pr.Title)

		reviewers := make([]*models.User, 0)
		for _, review := range reviews {
			user, err := transformUser(review.GetUser(), db)
			if err != nil {
				log.Println("Unable to transform reviewer user:", err)
				continue
			}

			reviewers = append(reviewers, user)

			if review.GetState() == approvedState {
				exists, err := user.ApprovedPullRequests(qm.Where("pull_request_id = ?", pr.ID)).Exists(context.Background(), db)
				if err != nil {
					log.Println("Unable to check pull request activity record:", err)
					continue
				}
				if !exists {
					err = user.AddApprovedPullRequests(context.Background(), db, false, pr)
					if err != nil {
						log.Println("Unable to persist user approved pull request")
						continue
					}
				}
			} else {
				exists, err := user.CommentedPullRequests(qm.Where("pull_request_id = ?", pr.ID)).Exists(context.Background(), db)
				if err != nil {
					log.Println("Unable to check pull request activity record:", err)
					continue
				}
				if !exists {
					err = user.AddCommentedPullRequests(context.Background(), db, false, pr)
					if err != nil {
						log.Println("Unable to persist user commented pull request")
						continue
					}
				}
			}
			if err != nil {
				log.Println("Unable to persist pull request activity record:", err)
			}
		}

		idlers := rf.findIdlers(pr.R.Reviewers, reviewers)
		pr.AddIdlers(context.Background(), db, false, idlers...)

		pr.Update(context.Background(), db, boil.Infer()) // updates the 'updated_at' column
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
