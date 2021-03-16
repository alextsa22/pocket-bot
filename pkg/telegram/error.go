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
	msg := tgbotapi.NewMessage(chatID, "unknown error")

	switch err {
	case errUnauthorized:
		msg.Text = "you are not authorized, use the /start command"
	case errInvalidURL:
		msg.Text = "this is an invalid link"
	case errUnableToSave:
		msg.Text = "failed to save link, please try again"
	}

	b.bot.Send(msg)
}
