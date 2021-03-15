package redisdb

import (
	"github.com/alextsa22/pocket-bot/pkg/repository"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"strconv"
)

type TokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) *TokenRepository {
	return &TokenRepository{client: client}
}

func (r *TokenRepository) Set(chatID int64, token string, tokenType repository.TokenType) error {
	_, err := r.client.HSet(string(tokenType), strconv.Itoa(int(chatID)), token).Result()

	return err
}

func (r *TokenRepository) Get(chatID int64, tokenType repository.TokenType) (string, error) {
	token, err := r.client.HGet(string(tokenType), strconv.Itoa(int(chatID))).Result()
	if err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token not found")
	}

	return token, nil
}
