package telegram

import (
	"context"
	"github.com/alextsa22/pocket-bot/pkg/pocket"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/url"
)

const (
	commandStart = "start"

	replyStartTemplate = "Hey! To save links in your Pocket account, first you need to give me access to it. To do this, follow the link:\n%s"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "you are authorized")
	b.bot.Send(msg)

	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "command not found")
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return errUnauthorized
	}

	_, err = url.Parse(message.Text)
	if err != nil {
		return errInvalidURL
	}

	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		return errUnableToSave
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "link saved successfully")
	_, err = b.bot.Send(msg)
	return err
}
