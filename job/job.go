package job

import (
	"database/sql"
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

// JobContainer constitutes a set of jobs that the application will execute
type JobContainer interface {
	Jobs() []Job
}

type defaultJobContainer struct {
	db *sql.DB
}

// DefaultContainer creates an instance of the default job container
func DefaultContainer(db *sql.DB) JobContainer {
	return defaultJobContainer{
		db: db,
	}
}

// Jobs returns the defined jobs for the default job container
func (jc defaultJobContainer) Jobs() []Job {
	return []Job{
		NewPullRequestFetcher(jc.db),
	}
}
