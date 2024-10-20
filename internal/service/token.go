package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/KarmaBeLike/crypto-service/internal/models"
	"github.com/KarmaBeLike/crypto-service/internal/repository"
)

const coingeckoAPI = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd"

type TokenService interface {
	FetchAndStoreTokens() error
}

type tokenService struct {
	tokenRepo repository.TokenRepository
}

func NewTokenService(tokenRepo repository.TokenRepository) TokenService {
	return &tokenService{tokenRepo: tokenRepo}
}

func (s *tokenService) FetchAndStoreTokens() error {
	tokens, err := s.fetchTokens()
	if err != nil {
		return err
	}

	// Сохраняем токены через репозиторий
	if err := s.tokenRepo.InsertTokens(tokens); err != nil {
		return err
	}

	return nil
}

func (s *tokenService) fetchTokens() ([]models.Token, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(coingeckoAPI)
	if err != nil {
		return nil, errors.New("failed to fetch tokens from coingecko: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code from coingecko: " + resp.Status)
	}

	var tokens []models.Token
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return tokens, nil
}
