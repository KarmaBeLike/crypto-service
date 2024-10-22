package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/KarmaBeLike/crypto-service/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
)

func OpenDB(cfg *config.Config) (*sql.DB, error) {
	dbURL := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to the database:")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "Database ping failed")
	}

	log.Println("Database connection established.")
	return db, nil
}

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "create driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return errors.Wrap(err, "migrate")
	}
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion { // если ошибка не ErrNilVersion (нет миграций)
		return errors.Wrap(err, "failed to get migration version")
	}

	if dirty {
		log.Printf("Dirty migration detected at version %d. Forcing to clean state.\n", version)

		if err := m.Force(int(version)); err != nil {
			return errors.Wrap(err, "failed to force migration")
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to apply migrations")
	}

	log.Println("Migrations applied successfully!")

	return nil
}
