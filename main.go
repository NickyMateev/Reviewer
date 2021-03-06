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

	storage, err := storage.New(config.Storage)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	srv := server.New(config.Server, storage, github.NewClient(nil), config.Slack)
	if err := srv.Run(); err != nil {
		panic(err)
	}
}
