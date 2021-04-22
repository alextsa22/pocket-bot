package main

import (
	"github.com/alextsa22/pocket-bot/internal/config"
	"github.com/alextsa22/pocket-bot/internal/repository/redisdb"
	"github.com/alextsa22/pocket-bot/internal/server"
	"github.com/alextsa22/pocket-bot/internal/telegram"
	"github.com/alextsa22/pocket-bot/pkg/pocket"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	telegramTokenEnv     = "TOKEN"
	pocketConsumerKeyEnv = "CONSUMER_KEY"
)

func main() {
	logrus.StandardLogger().SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.Init()
	if err != nil {
		logrus.WithError(err).Fatal("config init error")
	}
	logrus.Info("config initialized")

	bot, err := tgbotapi.NewBotAPI(os.Getenv(telegramTokenEnv))
	if err != nil {
		logrus.WithError(err).Fatal("bot init error")
	}
	logrus.Info("bot initialized")

	bot.Debug = false // true

	pocketClient, err := pocket.NewClient(os.Getenv(pocketConsumerKeyEnv))
	if err != nil {
		logrus.WithError(err).Fatal("pocket client init error")
	}
	logrus.Info("pocket client initialized")

	redisClient, err := initRedisClient()
	if err != nil {
		logrus.WithError(err).Fatal("redis client init error")
	}
	logrus.Info("redis client initialized")

	tokenRepo := redisdb.NewTokenRepository(redisClient)
	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepo, cfg.AuthServer.GetRedirectURL(), cfg.Messages)
	authorizationServer := server.NewAuthorizationServer(cfg.AuthServer.Port, pocketClient, tokenRepo, cfg.TelegramBotURL)

	go func() {
		if err := telegramBot.Start(); err != nil {
			logrus.WithError(err).Fatal("telegram bot start error")
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		logrus.WithError(err).Fatal("authorization server error")
	}
}

func initRedisClient() (*redis.Client, error) {
	redisCfg, err := config.InitRedisConfig()
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisCfg.Host + ":" + redisCfg.Port,
	})
	if _, err = client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}
