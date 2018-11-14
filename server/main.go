package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cheshir/tttnn/server/ai"
)

const (
	defaultTimeout = 15 * time.Second
)

func init() {
	// TODO configure logger.
	// TODO switch to go mod.
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", config)
	}

	log.Printf("config: %#v\n", config)

	predictor, err := ai.New(config.AI)
	if err != nil {
		log.Fatalf("Failed to initialize AI: %v", err)
	}

	defer predictor.Close()

	mux := http.NewServeMux()
	mux.Handle("/api/predict", newPredictHandler(predictor))

	server := &http.Server{
		Addr:              config.Host + ":" + config.Port,
		Handler:           mux,
		IdleTimeout:       defaultTimeout,
		ReadHeaderTimeout: defaultTimeout,
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
		ErrorLog:          log.New(os.Stdout, "server", log.LstdFlags|log.Lmicroseconds),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve HTTP API: %v", err)
	}
}
