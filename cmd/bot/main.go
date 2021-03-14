package main

import (
	"github.com/alextsa22/pocket-bot/pkg/pocket"
	"github.com/alextsa22/pocket-bot/pkg/telegram"
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
		log.Fatal("1", err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(os.Getenv("CONSUMER_KEY"))
	if err != nil {
		log.Fatal("2", err)
	}

	telegramBot := telegram.NewBot(bot, pocketClient, "http://localhost")
	if err := telegramBot.Start(); err != nil {
		log.Fatal("3", err)
	}
}
