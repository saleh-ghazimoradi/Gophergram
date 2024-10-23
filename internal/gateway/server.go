package gateway

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/docs"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"os"
	"os/signal"
	"syscall"

	"net/http"
	"time"
)

func Server(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = Version
	docs.SwaggerInfo.Host = config.AppConfig.General.APIURL.APIURLSwag
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         config.AppConfig.General.Listen,
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		logger.Logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdown
	if err != nil {
		return err
	}

	logger.Logger.Infow("server has stopped", "addr", config.AppConfig.General.Listen, "env", config.AppConfig.Env.Env)

	return nil
}
