package handlers

import (
	"net/http"

	"github.com/KarmaBeLike/crypto-service/internal/service"
)

// tokensHandler HTTP-обработчик для получения списка токенов и их цен
type TokenHandler struct {
	tokenService service.TokenService
}

func NewTokenHandler(tokenService service.TokenService) *TokenHandler {
	return &TokenHandler{tokenService: tokenService}
}

func (h *TokenHandler) GetAndStoreTokens(w http.ResponseWriter, r *http.Request) {
	if err := h.tokenService.FetchAndStoreTokens(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tokens successfully fetched and stored"))
}
