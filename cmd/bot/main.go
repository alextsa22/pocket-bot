package main

import (
	"github.com/alextsa22/pocket-bot/pkg/config"
	"github.com/alextsa22/pocket-bot/pkg/pocket"
	"github.com/alextsa22/pocket-bot/pkg/repository/redisdb"
	"github.com/alextsa22/pocket-bot/pkg/server"
	"github.com/alextsa22/pocket-bot/pkg/telegram"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

const (
	telegramTokenEnv     = "TOKEN"
	pocketConsumerKeyEnv = "CONSUMER_KEY"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv(telegramTokenEnv))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(os.Getenv(pocketConsumerKeyEnv))
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := initRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	tokenRepo := redisdb.NewTokenRepository(redisClient)
	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepo, cfg.AuthServer.GetRedirectURL(), cfg.Messages)
	authorizationServer := server.NewAuthorizationServer(cfg.AuthServer.Port, pocketClient, tokenRepo, cfg.TelegramBotURL)

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal()
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
