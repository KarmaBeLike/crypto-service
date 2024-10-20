package repository

import (
	"database/sql"
	"errors"

	"github.com/KarmaBeLike/crypto-service/internal/models"
)

type TokenRepository interface {
	InsertTokens(tokens []models.Token) error
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
		          ON CONFLICT (id) DO UPDATE SET price_usd = $3`
		_, err := r.db.Exec(query, token.Symbol, token.Name, token.PriceUSD)
		if err != nil {
			return errors.New("failed to insert token: " + err.Error())
		}
	}
	return nil
}
