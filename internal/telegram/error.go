package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

var (
	errInvalidURL   = errors.New("url is invalid")
	errUnauthorized = errors.New("user is not authorized")
	errUnableToSave = errors.New("unable to save")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, b.messages.Errors.Default)

	switch err {
	case errUnauthorized:
		msg.Text = b.messages.Errors.Unauthorized
	case errInvalidURL:
		msg.Text = b.messages.Errors.InvalidURL
	case errUnableToSave:
		msg.Text = b.messages.Errors.UnableToSave
	}

	b.bot.Send(msg)
}
