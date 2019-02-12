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
	db     *sql.DB
	client *github.Client
}

// NewPullRequestFetcher creates an instance of PullRequestFetcher
func NewPullRequestFetcher(db *sql.DB, client *github.Client) *PullRequestFetcher {
	return &PullRequestFetcher{
		db:     db,
		client: client,
	}
}

// Name returns the name of the PullRequestFetcher job
func (prf PullRequestFetcher) Name() string {
	return pullRequestFetcher
}

// Period returns the period of time when the PullRequestFetcher job should execute
func (prf PullRequestFetcher) Period() string {
	return "10s"
}

// Run executes the PullRequestFetcher job
func (prf PullRequestFetcher) Run() {
	log.Printf("STARTING %v job", prf.Name())

	projects, err := models.Projects().All(context.Background(), prf.db)
	if err != nil {
		log.Panic("Unable to fetch projects:", err)
	}
	log.Printf("%d project(s) are about to be reconciled\n", len(projects))

	for _, project := range projects {
		projectName := fmt.Sprintf("%q [%v/%v]", project.Name, project.RepoOwner, project.RepoName)
		log.Printf("Fetching pull requests for project %v\n", projectName)

		pullRequests, resp, err := prf.client.PullRequests.List(context.Background(), project.RepoOwner, project.RepoName, nil)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Panic("Unable to fetch pull requests:", err)
		}
		log.Printf("(%d) pull request(s) fetched for project %v\n", len(pullRequests), projectName)

		go prf.fetchPullRequests(pullRequests, projectName)
	}
}

func (prf PullRequestFetcher) fetchPullRequests(pullRequests []*github.PullRequest, projectName string) {
	for _, pullRequest := range pullRequests {
		exists, err := models.PullRequests(qm.Where("github_id = ?", pullRequest.GetID())).Exists(context.Background(), prf.db)
		if err != nil {
			log.Panicf("Error retrieving pull requests for %v: %v\n", projectName, err.Error())
		}

		if !exists {
			user := prf.transformUser(pullRequest.GetUser())

			log.Printf("Persisting new pull request: %q (%v)\n", pullRequest.GetTitle(), projectName)

			pr := models.PullRequest{
				Title:     pullRequest.GetTitle(),
				URL:       pullRequest.GetHTMLURL(),
				Number:    int64(pullRequest.GetNumber()),
				GithubID:  pullRequest.GetID(),
				UserID:    user.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now()}

			err = pr.Insert(context.Background(), prf.db, boil.Infer())
			if err != nil {
				log.Panicf("Error persisting pull request %q (%v): %v\n", pr.Title, projectName, err.Error())
			}
			log.Printf("Pull request %q successfully persisted (%v)\n", pr.Title, projectName)

			reviewers := make([]*models.User, 0)
			for _, reviewer := range pullRequest.RequestedReviewers {
				reviewers = append(reviewers, prf.transformUser(reviewer))
			}
			if len(reviewers) > 0 {
				err := pr.AddReviewers(context.Background(), prf.db, false, reviewers...)
				if err != nil {
					log.Panicf("Error persisting pull request reviewers %q (%v): %v\n", pr.Title, projectName, err.Error())
				}
			}
		}
	}
}

func (prf PullRequestFetcher) transformUser(githubUser *github.User) *models.User {
	exists, err := models.Users(qm.Where("github_id = ?", githubUser.GetID())).Exists(context.Background(), prf.db)
	if err != nil {
		log.Panicf("Error searching for user %q [%v]: %v\n", githubUser.GetLogin(), githubUser.GetID(), err.Error())
	}

	var user *models.User
	if !exists {
		user = &models.User{Username: githubUser.GetLogin(), GithubID: githubUser.GetID()}
		err := user.Insert(context.Background(), prf.db, boil.Infer())
		if err != nil {
			log.Panicf("Error persisting user %q [%v]: %v\n", githubUser.GetLogin(), githubUser.GetID(), err.Error())
		}
	} else {
		user, err = models.Users(qm.Where("github_id = ?", githubUser.GetID())).One(context.Background(), prf.db)
		if err != nil {
			log.Panicf("Error retrieving user %q [%v]: %v\n", githubUser.GetLogin(), githubUser.GetID(), err.Error())
		}
	}

	return user
}
