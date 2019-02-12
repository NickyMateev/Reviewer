package job

import (
	"database/sql"
	"github.com/google/go-github/github"
	"github.com/robfig/cron"
)

const (
	pullRequestFetcher = "PullRequestFetcher"
)

// Job defines a task that will be executed periodically
type Job interface {
	cron.Job
	Period() string
	Name() string
}

// Container constitutes a set of jobs that the application will execute
type Container interface {
	Jobs() []Job
}

type defaultContainer struct {
	db     *sql.DB
	client *github.Client
}

// DefaultContainer creates an instance of the default job container
func DefaultContainer(db *sql.DB, client *github.Client) Container {
	return defaultContainer{
		db:     db,
		client: client,
	}
}

// Jobs returns the defined jobs for the default job container
func (jc defaultContainer) Jobs() []Job {
	return []Job{
		NewPullRequestFetcher(jc.db, jc.client),
	}
}
