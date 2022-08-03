package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "development/production/tmp mode" )

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)

	app := application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.port),
		Handler:           app.routes(),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		IdleTimeout:      time.Minute,
	}

	logger.Printf("starting %s server on %d ", cfg.env, cfg.port)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalln(err)
	}
}
