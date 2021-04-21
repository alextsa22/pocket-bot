package repository

type TokenType string

const (
	AccessToken  TokenType = "access_tokens"
	RequestToken TokenType = "request_tokens"
)

type TokenRepository interface {
	Set(chatID int64, token string, tokenType TokenType) error
	Get(chatID int64, tokenType TokenType) (string, error)
}
