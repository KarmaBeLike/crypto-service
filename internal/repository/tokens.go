package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/KarmaBeLike/crypto-service/internal/models"
)

type TokenRepository interface {
	InsertTokens(tokens []models.Token) error
	GetTokenPriceHistory(tokenSymbol string) ([]models.TokenPriceHistory, error)
	InsertTokenPriceHistory(tokens []models.Token) error
}

type tokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) InsertTokens(tokens []models.Token) error {
	for _, token := range tokens {
		query := `INSERT INTO tokens (symbol, name, price_usd)
		          VALUES ($1, $2, $3)
		          ON CONFLICT (symbol) DO UPDATE SET price_usd = $3`
		_, err := r.db.Exec(query, token.Symbol, token.Name, token.PriceUSD)
		if err != nil {
			return errors.New("failed to insert token: " + err.Error())
		}
	}
	return nil
}

func (r *tokenRepository) InsertTokenPriceHistory(tokens []models.Token) error {
	for _, token := range tokens {
		var tokenID int
		err := r.db.QueryRow(`SELECT id FROM tokens WHERE symbol = $1`, token.Symbol).Scan(&tokenID)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("Token with symbol %s not found\n", token.Symbol)
				continue // Пропускаем, если токен не найден
			}
			return err
		}

		if token.PriceUSD == 0 {
			fmt.Println("Warning: Price is zero for token:", token.Symbol)
			continue
		}

		_, err = r.db.Exec(`
            INSERT INTO token_price_history (token_id, price, symbol)
            VALUES ($1, $2, $3)
        `, tokenID, token.PriceUSD, token.Symbol)
		if err != nil {
			return fmt.Errorf("failed to insert price history for token ID %d: %w", tokenID, err)
		}
	}
	return nil
}

func (r *tokenRepository) GetTokenPriceHistory(tokenID string) ([]models.TokenPriceHistory, error) {
	var history []models.TokenPriceHistory

	var dbTokenID int
	err := r.db.QueryRow(`SELECT id FROM tokens WHERE id = $1`, tokenID).Scan(&dbTokenID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no token found with the given id")
		}
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT price, created_at, symbol
		FROM token_price_history
		WHERE token_id = $1
		ORDER BY created_at DESC
	`, dbTokenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var historyEntry models.TokenPriceHistory
		if err := rows.Scan(&historyEntry.CurrentPrice, &historyEntry.CreatedAt, &historyEntry.Symbol); err != nil {
			return nil, err
		}
		history = append(history, historyEntry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}
