package main

import (
	"fmt"
	"github.com/NickyMateev/Reviewer/config"
	"github.com/NickyMateev/Reviewer/server"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/google/go-github/github"
	"os"
)

func main() {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	err = setCFEnvVariables(config)
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

func setCFEnvVariables(config *config.Settings) error {
	if os.Getenv("VCAP_APPLICATION") != "" {
		appEnv, err := cfenv.Current()
		if err != nil {
			return fmt.Errorf("could not load CF environment: %v", err)
		}

		config.Server.Port = appEnv.Port

		storageInstance := os.Getenv("STORAGE_NAME")
		if storageInstance == "" {
			return fmt.Errorf("could not find storage instance name in environment")
		}

		service, err := appEnv.Services.WithName(storageInstance)
		if err != nil {
			return fmt.Errorf("could not retrieve storage information from environment")
		}
		config.Storage.URI = service.Credentials["uri"].(string) + "?sslmode=disable"

		config.Storage.MigrationsURL = os.Getenv("STORAGE_MIGRATIONS_URL")
	}
	return nil
}
