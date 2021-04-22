package server

import (
	"github.com/alextsa22/pocket-bot/internal/repository"
	"github.com/alextsa22/pocket-bot/pkg/pocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type AuthorizationServer struct {
	port         string
	server       *http.Server
	pocketClient *pocket.Client
	tokenRepo    repository.TokenRepository
	redirectURL  string
}

func NewAuthorizationServer(port string, pocketClient *pocket.Client, tokenRepo repository.TokenRepository, redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{port: port, pocketClient: pocketClient, tokenRepo: tokenRepo, redirectURL: redirectURL}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s,
	}

	logrus.Info("authorization server is running")
	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
		}).Error("invalid method")

		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIDParam := r.URL.Query().Get("chat_id")
	if chatIDParam == "" {
		logrus.WithField("url", r.URL.String()).Error("no chat_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		logrus.WithField("chatIdParam", chatIDParam).
			WithError(err).Error("invalid chat_id type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := s.tokenRepo.Get(chatID, repository.RequestToken)
	if err != nil {
		logrus.WithField("chatId", chatID).
			WithError(err).Error("request token not found")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authResp, err := s.pocketClient.Authorize(r.Context(), requestToken)
	if err != nil {
		logrus.WithError(err).Error("error while getting access token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.tokenRepo.Set(chatID, authResp.AccessToken, repository.AccessToken)
	if err != nil {
		logrus.WithError(err).Error("access token saving error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
