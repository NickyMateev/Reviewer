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
	"sync"
	"time"
)

// PullRequestFetcher is a regular job which fetches the pull requests for all registered projects
type PullRequestFetcher struct {
	storage storage.Storage
	client  *github.Client
}

// NewPullRequestFetcher creates an instance of PullRequestFetcher
func NewPullRequestFetcher(storage storage.Storage, client *github.Client) *PullRequestFetcher {
	return &PullRequestFetcher{
		storage: storage,
		client:  client,
	}
}

// Name returns the name of the PullRequestFetcher job
func (prf *PullRequestFetcher) Name() string {
	return pullRequestFetcher
}

// Period returns the period of time when the PullRequestFetcher job should execute
func (prf *PullRequestFetcher) Period() string {
	return "@every 1h30m"
}

// Run executes the PullRequestFetcher job
func (prf *PullRequestFetcher) Run() {
	log.Printf("STARTING %v job", prf.Name())
	defer log.Printf("FINISHED %v job", prf.Name())

	projects, err := models.Projects().All(context.Background(), prf.storage.Get())
	if err != nil {
		log.Println("Unable to fetch projects:", err)
		return
	}
	log.Printf("%d project(s) are about to be reconciled\n", len(projects))

	wg := sync.WaitGroup{}
	for _, project := range projects {
		projectName := fmt.Sprintf("%q [%v/%v]", project.Name, project.RepoOwner, project.RepoName)
		log.Printf("Fetching pull requests for project %v\n", projectName)

		pullRequests, resp, err := prf.client.PullRequests.List(context.Background(), project.RepoOwner, project.RepoName, nil)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Println("Unable to fetch pull requests:", err)
			continue
		}
		log.Printf("(%d) pull request(s) fetched for project %v\n", len(pullRequests), projectName)

		wg.Add(1)
		go prf.persistPullRequests(pullRequests, project.ID, projectName, &wg)
	}
	wg.Wait()
}

func (prf *PullRequestFetcher) persistPullRequests(pullRequests []*github.PullRequest, projectID int64, projectName string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, pullRequest := range pullRequests {
		txErr := prf.storage.Transaction(context.Background(), func(context context.Context, tx *sql.Tx) error {
			exists, err := models.PullRequests(qm.Where("github_id = ?", pullRequest.GetID())).Exists(context, tx)
			if err != nil {
				return fmt.Errorf("error retrieving pull requests for %v: %s", projectName, err)
			}

			if exists {
				return nil
			}

			user, err := transformUser(context, tx, pullRequest.GetUser())
			if err != nil {
				return fmt.Errorf("unable to transform user: %s", err)
			}

			log.Printf("Persisting new pull request: %q (%v)\n", pullRequest.GetTitle(), projectName)

			pr := models.PullRequest{
				Title:     pullRequest.GetTitle(),
				URL:       pullRequest.GetHTMLURL(),
				Number:    int64(pullRequest.GetNumber()),
				GithubID:  pullRequest.GetID(),
				UserID:    user.ID,
				ProjectID: projectID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now()}

			err = pr.Insert(context, tx, boil.Infer())
			if err != nil {
				return err
			}
			log.Printf("Pull request %q successfully persisted (%v)\n", pr.Title, projectName)

			reviewers := make([]*models.User, 0)
			for _, reviewer := range pullRequest.RequestedReviewers {
				rev, err := transformUser(context, tx, reviewer)
				if err != nil {
					return fmt.Errorf("unable to transform reviewer user: %s", err)
				}

				reviewers = append(reviewers, rev)
			}
			if len(reviewers) > 0 {
				err := pr.AddReviewers(context, tx, false, reviewers...)
				if err != nil {
					return err
				}
			}
			return nil
		})

		if txErr != nil {
			log.Printf("Unable to persist pull request %q: %s\n", *pullRequest.Title, txErr)
		}
	}
}
