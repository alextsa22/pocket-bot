package main

import (
	"github.com/alextsa22/pocket-bot/pkg/pocket"
	"github.com/alextsa22/pocket-bot/pkg/repository/redisdb"
	"github.com/alextsa22/pocket-bot/pkg/server"
	"github.com/alextsa22/pocket-bot/pkg/telegram"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(os.Getenv("CONSUMER_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if _, err := redisClient.Ping().Result(); err != nil {
		log.Fatal(err)
	}

	tokenRepo := redisdb.NewTokenRepository(redisClient)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepo, "http://localhost/")

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepo, "https://t.me/pocket_x_bot")

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal()
	}
}
