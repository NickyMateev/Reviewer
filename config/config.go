package config

import (
	"github.com/NickyMateev/Reviewer/server"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/spf13/viper"
)

// Settings consists of all the application configuration
type Settings struct {
	Storage storage.Config
	Server  server.Config
}

// New creates an instance of the application's configuration
func New() (*Settings, error) {
	v := viper.New()

	v.SetConfigName("application")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Settings{
		Storage: storage.Config{
			Type: v.GetString("storage.type"),
			URI:  v.GetString("storage.uri"),
		},
		Server: server.Config{
			Port:           v.GetInt("server.port"),
			RequestTimeout: v.GetDuration("server.request_timeout"),
		},
	}, nil
}
