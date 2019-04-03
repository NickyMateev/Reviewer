package job

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/google/go-github/github"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

type SlackConfig struct {
	ChannelID string
	BotToken  string
}

func transformUser(context context.Context, tx *sql.Tx, githubUser *github.User) (*models.User, error) {
	exists, err := models.Users(qm.Where("github_id = ?", githubUser.GetID())).Exists(context, tx)
	if err != nil {
		return nil, fmt.Errorf("Error searching for user %q [%v]: %v\n", githubUser.GetLogin(), githubUser.GetID(), err.Error())
	}

	var user *models.User
	if !exists {
		user = &models.User{Username: githubUser.GetLogin(), GithubID: githubUser.GetID()}
		err := user.Insert(context, tx, boil.Infer())
		if err != nil {
			return nil, fmt.Errorf("Error persisting user %q [%v]: %v\n", githubUser.GetLogin(), githubUser.GetID(), err.Error())
		}
	} else {
		user, err = models.Users(qm.Where("github_id = ?", githubUser.GetID())).One(context, tx)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving user %q [%v]: %v\n", githubUser.GetLogin(), githubUser.GetID(), err.Error())
		}
	}

	return user, nil
}
