package main

import (
	"github.com/NickyMateev/Reviewer/config"
	"github.com/NickyMateev/Reviewer/server"
	"github.com/NickyMateev/Reviewer/storage"
)

func main() {
	config, err := config.New()

	db, err := storage.New(config.Storage)
	if err != nil {
		panic(err)
	}

	srv, err := server.New(config.Server, db)
	if err != nil {
		panic(err)
	}

	srv.Run()
}
