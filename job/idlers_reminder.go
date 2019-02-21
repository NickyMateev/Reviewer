package job

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/nlopes/slack"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"log"
)

const githubURL = "https://github.com/%v/%v"

// IdlersReminder is a regular job which sends a notification to all users who have not reviewed assigned pull requests
type IdlersReminder struct {
	db     *sql.DB
	client *slack.Client
	config SlackConfig
}

// NewIdlersReminder creates an instance of IdlersReminder
func NewIdlersReminder(db *sql.DB, config SlackConfig) *IdlersReminder {
	return &IdlersReminder{
		db:     db,
		client: slack.New(config.BotToken),
		config: config,
	}
}

// Name returns the name of the IdlersReminder job
func (ir *IdlersReminder) Name() string {
	return idlersReminder
}

// Period returns the period of time when the IdlersReminder job should execute
func (ir *IdlersReminder) Period() string {
//	return "0 0 16 * * *"
	return "@every 1m"
}

// Run executes the IdlersReminder job
func (ir *IdlersReminder) Run() {
	log.Printf("STARTING %v job", ir.Name())
	defer log.Printf("FINISHED %v job", ir.Name())

	pullRequests, err := models.PullRequests(qm.Load("Idlers"), qm.Load("Project")).All(context.Background(), ir.db)
	if err != nil {
		log.Panic("Error retrieving pull requests:", err)
	}

	attachment := new(slack.Attachment)
	for _, pullRequest := range pullRequests {
		for _, idler := range pullRequest.R.Idlers {
			userMetadata := struct {
				SlackID string `json:"slack_id"`
			}{
				SlackID: ir.config.ChannelID,
			}

			err := idler.Metadata.Unmarshal(&userMetadata)
			if err != nil {
				log.Printf("Error unmarshalling user [%v] metadata:%v\n", idler.Username, err)
			}

			attachment.Title = fmt.Sprintf("[%v]", pullRequest.R.Project.Name)
			attachment.TitleLink = fmt.Sprintf(githubURL, pullRequest.R.Project.RepoOwner, pullRequest.R.Project.RepoName)
			attachment.Text = fmt.Sprintf(":arrow_right: *%v [#%v]*\n\t%v", pullRequest.Title, pullRequest.Number, pullRequest.URL)

			var mentionedUser string
			if userMetadata.SlackID != ir.config.ChannelID {
				mentionedUser = fmt.Sprintf("<@%v>", userMetadata.SlackID)
			} else {
				mentionedUser = fmt.Sprintf("%v", idler.Username)
			}
			_, _, err = ir.client.PostMessage(userMetadata.SlackID,
				slack.MsgOptionText(fmt.Sprintf("*[%v]* :rotating_light: Review the following pull request unless you want to be part of the *Wall of Shame* tomorrow :rotating_light:", mentionedUser), false),
				slack.MsgOptionAttachments(*attachment))
			if err != nil {
				log.Println("Slack notification could not be sent:", err)
			} else {
				log.Printf("Slack notification sent to %v for pending Pull Request review on %q\n", idler.Username, pullRequest.Title)
			}

		}
	}
}
