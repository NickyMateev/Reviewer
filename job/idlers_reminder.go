package job

import (
	"context"
	"fmt"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/nlopes/slack"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"log"
)

// IdlersReminder is a regular job which sends a notification to all users who have not reviewed assigned pull requests
type IdlersReminder struct {
	storage storage.Storage
	client  *slack.Client
	config  SlackConfig
}

// NewIdlersReminder creates an instance of IdlersReminder
func NewIdlersReminder(storage storage.Storage, config SlackConfig) *IdlersReminder {
	return &IdlersReminder{
		storage: storage,
		client:  slack.New(config.BotToken),
		config:  config,
	}
}

// Name returns the name of the IdlersReminder job
func (ir *IdlersReminder) Name() string {
	return idlersReminder
}

// Period returns the period of time when the IdlersReminder job should execute
func (ir *IdlersReminder) Period() string {
	return "0 0 16 * * *"
}

// Run executes the IdlersReminder job
func (ir *IdlersReminder) Run() {
	log.Printf("STARTING %v job", ir.Name())
	defer log.Printf("FINISHED %v job", ir.Name())

	pullRequests, err := models.PullRequests(qm.Load("Idlers"), qm.Load("Project")).All(context.Background(), ir.storage.Get())
	if err != nil {
		log.Panic("Error retrieving pull requests:", err)
	}

	for _, pullRequest := range pullRequests {
		for _, idler := range pullRequest.R.Idlers {
			msg := ir.prepareMessage(idler, pullRequest, pullRequest.R.Project)
			if err := ir.sendSlackNotification(msg); err != nil {
				log.Println("Slack notification could not be sent:", err)
			} else {
				log.Printf("Slack notification sent to %v for pending Pull Request review on %q\n", idler.Username, pullRequest.Title)
			}
		}
	}
}

func (ir IdlersReminder) prepareMessage(idler *models.User, pullRequest *models.PullRequest, project *models.Project) *Message {
	userMetadata := struct {
		SlackID string `json:"slack_id"`
	}{
		SlackID: ir.config.ChannelID,
	}

	err := idler.Metadata.Unmarshal(&userMetadata)
	if err != nil {
		log.Printf("Error unmarshalling user [%v] metadata:%v\n", idler.Username, err)
	}

	var mentionedUser string
	if userMetadata.SlackID != ir.config.ChannelID {
		mentionedUser = fmt.Sprintf("<@%v>", userMetadata.SlackID)
	} else {
		mentionedUser = fmt.Sprintf("%v", idler.Username)
	}
	header := fmt.Sprintf("*[%v]* :rotating_light: Review the following pull request unless you want to be part of the *Wall of Shame* tomorrow :rotating_light:", mentionedUser)

	attachment := new(slack.Attachment)
	attachment.Title = fmt.Sprintf("[%v]", project.Name)
	attachment.TitleLink = extractURLFromProject(project)

	attachment.Text = fmt.Sprintf(":arrow_right: *%v [#%v]*\n\t%v", pullRequest.Title, pullRequest.Number, pullRequest.URL)

	return &Message{
		Attachment: attachment,
		Header:     header,
		Receiver:   userMetadata.SlackID,
	}
}

func (ir IdlersReminder) sendSlackNotification(msg *Message) error {
	_, _, err := ir.client.PostMessage(msg.Receiver,
		slack.MsgOptionText(fmt.Sprintf(msg.Header), false),
		slack.MsgOptionAttachments(*msg.Attachment))
	return err
}
