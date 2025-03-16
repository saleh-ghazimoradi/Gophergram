package utils

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
)

func PostURI() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.AppConfig.Database.DatabaseHost, config.AppConfig.Database.DatabasePort, config.AppConfig.Database.DatabaseUser, config.AppConfig.Database.DatabasePassword, config.AppConfig.Database.DatabaseName, config.AppConfig.Database.DatabaseSSLMode)
}

func PostConnection() (*sql.DB, error) {
	postURI := PostURI()
	sLogger.SLogger.Info("connecting to postgres with options: " + postURI)

	db, err := sql.Open("postgres", postURI)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Postgres: %v", err)
	}
	db.SetMaxOpenConns(config.AppConfig.Database.MaxOpenConn)
	db.SetMaxIdleConns(config.AppConfig.Database.MaxIdleConn)
	db.SetConnMaxLifetime(config.AppConfig.Database.MaxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig.Database.Timeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging Postgres database: %w", err)
	}

	return db, nil
}
