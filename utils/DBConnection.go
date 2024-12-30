package utils

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
)

func PostURI() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.AppConfig.DBConfig.DbHost, config.AppConfig.DBConfig.DbPort, config.AppConfig.DBConfig.DbUser, config.AppConfig.DBConfig.DbPassword, config.AppConfig.DBConfig.DbName, config.AppConfig.DBConfig.DbSslMode)
}

func PostConnection() (*sql.DB, error) {
	postURI := PostURI()
	logger.Logger.Info("Connecting to Postgres with options: " + postURI)

	db, err := sql.Open("postgres", postURI)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Postgres: %v", err)
	}
	db.SetMaxOpenConns(config.AppConfig.DBConfig.MaxOpenConns)
	db.SetMaxIdleConns(config.AppConfig.DBConfig.MaxIdleConns)
	db.SetConnMaxLifetime(config.AppConfig.DBConfig.MaxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig.DBConfig.Timeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging Postgres database: %w", err)
	}

	return db, nil
}
