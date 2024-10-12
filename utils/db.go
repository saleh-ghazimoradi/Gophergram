package utils

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"time"

	_ "github.com/lib/pq" // Postgres driver
)

type PostgresConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Database     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
	Timeout      time.Duration
}

func PostgresURI(cfg PostgresConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func PostgresConnection(cfg PostgresConfig) (*sql.DB, error) {
	connString := PostgresURI(cfg)

	logger.Logger.Info("Connecting to Postgres with options: " + connString)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("error opening Postgres connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging Postgres database: %w", err)
	}

	return db, nil
}
