package job

import (
	"context"
	"database/sql"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/google/go-github/github"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"log"
)

type SlackConfig struct {
	ChannelID string
	BotToken  string
}

func transformUser(githubUser *github.User, db *sql.DB) *models.User {
	exists, err := models.Users(qm.Where("github_id = ?", githubUser.GetID())).Exists(context.Background(), db)
	if err != nil {
		log.Panicf("Error searching for user %q [%v]: %v\n", githubUser.GetLogin(), githubUser.GetID(), err.Error())
	}

	var user *models.User
	if !exists {
		user = &models.User{Username: githubUser.GetLogin(), GithubID: githubUser.GetID()}
		err := user.Insert(context.Background(), db, boil.Infer())
		if err != nil {
			log.Panicf("Error persisting user %q [%v]: %v\n", githubUser.GetLogin(), githubUser.GetID(), err.Error())
		}
	} else {
		user, err = models.Users(qm.Where("github_id = ?", githubUser.GetID())).One(context.Background(), db)
		if err != nil {
			log.Panicf("Error retrieving user %q [%v]: %v\n", githubUser.GetLogin(), githubUser.GetID(), err.Error())
		}
	}

	return user
}
