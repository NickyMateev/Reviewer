package job

import (
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/google/go-github/github"
	"github.com/robfig/cron"
)

const (
	pullRequestFetcher = "PullRequestFetcher"
	reviewFetcher      = "ReviewFetcher"
	idlersReminder     = "IdlersReminder"
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
	storage     storage.Storage
	client      *github.Client
	slackConfig SlackConfig
}

// DefaultContainer creates an instance of the default job container
func DefaultContainer(storage storage.Storage, client *github.Client, config SlackConfig) Container {
	return defaultContainer{
		storage:     storage,
		client:      client,
		slackConfig: config,
	}
}

// Jobs returns the defined jobs for the default job container
func (jc defaultContainer) Jobs() []Job {
	return []Job{
		NewPullRequestFetcher(jc.storage, jc.client),
		NewReviewFetcher(jc.storage, jc.client),
		NewIdlersReminder(jc.storage, jc.slackConfig),
	}
}
