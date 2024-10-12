package gateway

import (
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"log"
	"net/http"
	"time"
)

func Server() {
	srv := &http.Server{
		Addr:         config.AppConfig.General.Listen,
		Handler:      routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
