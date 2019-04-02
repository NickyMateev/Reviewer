package main

import (
	"github.com/NickyMateev/Reviewer/config"
	"github.com/NickyMateev/Reviewer/server"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/google/go-github/github"
)

func main() {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	db, err := storage.New(config.Storage)
	if err != nil {
		panic(err)
	}

	srv, err := server.New(config.Server, db, github.NewClient(nil), config.Slack)
	if err != nil {
		panic(err)
	}

	srv.Run()
}
