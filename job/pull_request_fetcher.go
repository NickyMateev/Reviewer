package job

import (
	"context"
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
		go prf.fetchPullRequests(pullRequests, project.ID, projectName, &wg)
	}
	wg.Wait()
}

func (prf *PullRequestFetcher) fetchPullRequests(pullRequests []*github.PullRequest, projectID int64, projectName string, wg *sync.WaitGroup) {
	defer wg.Done()
	db := prf.storage.Get()
	for _, pullRequest := range pullRequests {
		exists, err := models.PullRequests(qm.Where("github_id = ?", pullRequest.GetID())).Exists(context.Background(), db)
		if err != nil {
			log.Printf("Error retrieving pull requests for %v: %v\n", projectName, err.Error())
			continue
		}

		if exists {
			continue
		}

		user, err := transformUser(pullRequest.GetUser(), db)
		if err != nil {
			log.Println("Unable to transform user:", err)
			continue
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

		err = pr.Insert(context.Background(), db, boil.Infer())
		if err != nil {
			log.Printf("Error persisting pull request %q (%v): %v\n", pr.Title, projectName, err.Error())
			continue
		}
		log.Printf("Pull request %q successfully persisted (%v)\n", pr.Title, projectName)

		reviewers := make([]*models.User, 0)
		for _, reviewer := range pullRequest.RequestedReviewers {
			rev, err := transformUser(reviewer, db)
			if err != nil {
				log.Println("Unable to transform reviewer user:", err)
				continue
			}

			reviewers = append(reviewers, rev)
		}
		if len(reviewers) > 0 {
			err := pr.AddReviewers(context.Background(), db, false, reviewers...)
			if err != nil {
				log.Printf("Error persisting pull request reviewers %q (%v): %v\n", pr.Title, projectName, err.Error())
			}
		}
	}
}
