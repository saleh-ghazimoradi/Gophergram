package gateway

import (
	"context"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/docs"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/routes"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"github.com/saleh-ghazimoradi/Gophergram/utils"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

func Server() error {
	docs.SwaggerInfo.Version = config.AppConfig.ServerConfig.Version
	docs.SwaggerInfo.Host = config.AppConfig.ServerConfig.APIURL

	db, err := utils.PostConnection()
	if err != nil {
		return err
	}

	redis, err := utils.RedisConnection(config.AppConfig.Redis.Addr, config.AppConfig.Redis.PW, config.AppConfig.Redis.DB)
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	router := httprouter.New()
	routes.RegisterRoutes(router, db, redis)

	srv := &http.Server{
		Addr:         config.AppConfig.ServerConfig.Port,
		Handler:      router,
		ReadTimeout:  config.AppConfig.ServerConfig.ReadTimeout,
		WriteTimeout: config.AppConfig.ServerConfig.WriteTimeout,
		IdleTimeout:  config.AppConfig.ServerConfig.IdleTimeout,
	}

	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.Logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		logger.Logger.Info("completing background tasks", "addr", srv.Addr)

		wg.Wait()
		shutdownError <- nil
	}()

	logger.Logger.Info("starting server", "addr", config.AppConfig.ServerConfig.Port, "env", config.AppConfig.ServerConfig.Env)

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	logger.Logger.Info("stopped server", "addr", srv.Addr)

	return nil
}
