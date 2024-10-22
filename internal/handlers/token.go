package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KarmaBeLike/crypto-service/internal/service"
	"github.com/gorilla/mux"
)

// TokenHandler HTTP-обработчик для получения списка токенов и их цен
type TokenHandler struct {
	tokenService service.TokenService
}

func NewTokenHandler(tokenService service.TokenService) *TokenHandler {
	return &TokenHandler{tokenService: tokenService}
}

func (h *TokenHandler) GetAndStoreTokens(w http.ResponseWriter, r *http.Request) {
	tokens, err := h.tokenService.FetchAndStoreTokens()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		http.Error(w, "Failed to encode tokens to JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TokenHandler) GetTokenPriceHistory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler GetTokenPriceHistory called")
	vars := mux.Vars(r)
	tokenID := vars["token_id"]

	if tokenID == "" {
		http.Error(w, "token_id is missing", http.StatusBadRequest)
		return
	}

	fmt.Println("Received token ID:", tokenID)

	history, err := h.tokenService.GetTokenPriceHistory(tokenID)
	if err != nil {
		http.Error(w, "Failed to fetch price history: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверка на наличие истории
	if len(history) == 0 {
		http.Error(w, "No price history found for token", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(history); err != nil {
		http.Error(w, "Failed to encode history to JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
