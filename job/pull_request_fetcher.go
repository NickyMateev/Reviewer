package job

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/google/go-github/github"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"log"
	"net/http"
	"time"
)

// PullRequestFetcher is a regular job which fetches the pull requests for all registered projects
type PullRequestFetcher struct {
	db *sql.DB
}

// NewPullRequestFetcher creates an instance of PullRequestFetcher
func NewPullRequestFetcher(db *sql.DB) *PullRequestFetcher {
	return &PullRequestFetcher{
		db: db,
	}
}

// Name returns the name of the PullRequestFetcher job
func (prf PullRequestFetcher) Name() string {
	return pullRequestFetcher
}

// Period returns the period of time when the PullRequestFetcher job should execute
func (prf PullRequestFetcher) Period() string {
	return "30m"
}

// Run executes the PullRequestFetcher job
func (prf PullRequestFetcher) Run() {
	log.Printf("STARTING %v job", prf.Name())

	projects, err := models.Projects().All(context.Background(), prf.db)
	if err != nil {
		log.Panic("Unable to fetch projects:", err)
	}
	log.Printf("%d project(s) are about to be reconciled\n", len(projects))

	client := github.NewClient(nil)
	for _, project := range projects {
		projectName := fmt.Sprintf("%q [%v/%v]", project.Name, project.RepoOwner, project.RepoName)
		log.Printf("Fetching pull requests for project %v\n", projectName)

		pullRequests, resp, err := client.PullRequests.List(context.Background(), project.RepoOwner, project.RepoName, nil)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Panic("Unable to fetch pull requests:", err)
		}
		log.Printf("(%d) pull request(s) fetched for project %v\n", len(pullRequests), projectName)

		go prf.fetchPullRequests(pullRequests, projectName)
	}
}

func (prf PullRequestFetcher) fetchPullRequests(pullRequests []*github.PullRequest, projectName string) {
	for _, pr := range pullRequests {
		exists, err := models.PullRequests(qm.Where("github_id = ?", pr.GetID())).Exists(context.Background(), prf.db)
		if err != nil {
			log.Panicf("Error retrieving pull requests for %v: %v\n", projectName, err.Error())
		}

		if !exists {
			log.Printf("Persisting new pull request: %q (%v)\n", pr.GetTitle(), projectName)
			pr := models.PullRequest{Title: pr.GetTitle(), URL: pr.GetHTMLURL(), GithubID: pr.GetID(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
			err := pr.Insert(context.Background(), prf.db, boil.Infer())
			if err != nil {
				log.Panicf("Pull request %q could not be persisted: %v (%v)\n", pr.Title, err.Error(), projectName)
			}
			log.Printf("Pull request %q successfully persisted (%v)\n", pr.Title, projectName)
		}
	}

}
