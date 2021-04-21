package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramBotURL string `mapstructure:"bot_url"`
	AuthServer     AuthServer
	Messages       Messages
}

type AuthServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func (s *AuthServer) GetRedirectURL() string {
	return "http://" + s.Host + ":" + s.Port
}

type Messages struct {
	Responses Responses
	Errors    Errors
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	SavedSuccessfully string `mapstructure:"saved_successfully"`
	UnknownCommand    string `mapstructure:"unknown_command"`
}

type Errors struct {
	Default      string `mapstructure:"default"`
	Unauthorized string `mapstructure:"unauthorized"`
	InvalidURL   string `mapstructure:"invalid_url"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("auth_server", &cfg.AuthServer); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	return &cfg, nil
}
