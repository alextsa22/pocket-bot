package telegram

import (
	"context"
	"github.com/alextsa22/pocket-bot/pkg/pocket"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"net/url"
)

const (
	commandStart = "start"
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

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.AlreadyAuthorized)
	b.bot.Send(msg)

	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.UnknownCommand)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		logrus.WithError(err).Error(errUnauthorized)
		return errUnauthorized
	}

	if !isValidUrl(message.Text) {
		logrus.WithField("url", message.Text).Error(errInvalidURL)
		return errInvalidURL
	}

	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		logrus.WithError(err).Error(errUnableToSave)
		return errUnableToSave
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.SavedSuccessfully)
	_, err = b.bot.Send(msg)
	return err
}

func isValidUrl(text string) bool {
	_, err := url.ParseRequestURI(text)
	if err != nil {
		return false
	}

	u, err := url.Parse(text)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
