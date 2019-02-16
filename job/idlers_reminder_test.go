package job

import (
	"fmt"
	"github.com/NickyMateev/Reviewer/models"
	"testing"
)

var testIdlersJob = IdlersReminder{
	config: SlackConfig{
		ChannelID: "CG123",
	},
}

func TestPrepareMessageWhenUserHasSlackIDTagsPerson(t *testing.T) {
	u := &models.User{Username: "username", Metadata: []byte("{\"slack_id\":\"UG123\"}")}
	pr := &models.PullRequest{
		Title:  "PR Title",
		Number: 1,
		URL:    "http://domain.com",
	}
	proj := &models.Project{
		Name:      "Test Project",
		RepoOwner: "Project Owner",
		RepoName:  "Project Name",
	}

	msg := testIdlersJob.prepareMessage(u, pr, proj)

	expectedHeader := fmt.Sprint("*[<@UG123>]* :rotating_light: Review the following pull request unless you want to be part of the *Wall of Shame* tomorrow :rotating_light:")
	if msg.Header != expectedHeader {
		t.Errorf("header mismatch\nactual:%v\nexpected:%v", msg.Header, expectedHeader)
	}

	verifyMessageAttachment(t, msg, pr, proj)
}

func TestPrepareMessageWhenUserDoesntHaveSlackIDDoesntTagPerson(t *testing.T) {
	u := &models.User{Username: "username"}
	pr := &models.PullRequest{
		Title:  "PR Title",
		Number: 1,
		URL:    "http://domain.com",
	}
	proj := &models.Project{
		Name:      "Test Project",
		RepoOwner: "Project Owner",
		RepoName:  "Project Name",
	}

	msg := testIdlersJob.prepareMessage(u, pr, proj)

	expectedHeader := fmt.Sprintf("*[%v]* :rotating_light: Review the following pull request unless you want to be part of the *Wall of Shame* tomorrow :rotating_light:", u.Username)
	if msg.Header != expectedHeader {
		t.Errorf("header mismatch\nactual:%v\nexpected:%v", msg.Header, expectedHeader)
	}

	verifyMessageAttachment(t, msg, pr, proj)
}

func verifyMessageAttachment(t *testing.T, msg *Message, pr *models.PullRequest, proj *models.Project) {
	expectedTitle := "[" + proj.Name + "]"
	if msg.Title != expectedTitle {
		t.Errorf("title mismatch\nactual:%v\nexpected:%v", msg.Title, expectedTitle)
	}

	expectedTitleLink := extractURLFromProject(proj)
	if msg.TitleLink != expectedTitleLink {
		t.Errorf("title link mismatch\nactual:%v\nexpected:%v", msg.TitleLink, expectedTitleLink)
	}

	expectedText := fmt.Sprintf(":arrow_right: *%v [#%v]*\n\t%v", pr.Title, pr.Number, pr.URL)
	if msg.Text != expectedText {
		t.Errorf("text mismatch\nactual:%v\nexpected:%v", msg.Text, expectedText)
	}
}
