package gateway

import (
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"log"
	"net/http"
	"time"
)

func Server(mux http.Handler) error {
	srv := &http.Server{
		Addr:         config.AppConfig.General.Listen,
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	return nil
}
