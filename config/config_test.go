package config

import (
	"github.com/NickyMateev/Reviewer/job"
	"github.com/NickyMateev/Reviewer/server"
	"github.com/NickyMateev/Reviewer/storage"
	"testing"
	"time"
)

var expectedConfiguration = Settings{
	Server: server.Config{
		Port:           8888,
		RequestTimeout: 4000 * time.Millisecond,
	},
	Storage: storage.Config{
		Type: "postgres",
		URI:  "postgres://user:password@domain:5432/postgres",
	},
	Slack: job.SlackConfig{
		ChannelID: "CG1234567",
		BotToken:  "abc",
	},
}

func TestNewConfiguration(t *testing.T) {
	config, err := New()
	if err != nil {
		t.Errorf("unexpected error occurred while reading configuration into viper")
	}

	if config.Server.Port != expectedConfiguration.Server.Port {
		t.Errorf("server.port mismatch; actual=%v; expected=%v;\n", config.Server.Port, expectedConfiguration.Server.Port)
	}

	if config.Server.RequestTimeout != expectedConfiguration.Server.RequestTimeout {
		t.Errorf("server.request_timeout mismatch; actual=%v; expected=%v;\n", config.Server.RequestTimeout, expectedConfiguration.Server.RequestTimeout)
	}

	if config.Storage.Type != expectedConfiguration.Storage.Type {
		t.Errorf("storage.type mismatch; actual=%v; expected=%v;\n", config.Storage.Type, expectedConfiguration.Storage.Type)
	}

	if config.Storage.URI != expectedConfiguration.Storage.URI {
		t.Errorf("storage.uri mismatch; actual=%v; expected=%v;\n", config.Storage.URI, expectedConfiguration.Storage.URI)
	}

	if config.Slack.ChannelID != expectedConfiguration.Slack.ChannelID {
		t.Errorf("slack.channel_id mismatch; actual=%v; expected=%v;\n", config.Slack.ChannelID, expectedConfiguration.Slack.ChannelID)
	}

	if config.Slack.BotToken != expectedConfiguration.Slack.BotToken {
		t.Errorf("slack.bot_token mismatch; actual=%v; expected=%v;\n", config.Slack.BotToken, expectedConfiguration.Slack.BotToken)
	}
}
