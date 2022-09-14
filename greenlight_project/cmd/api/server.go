package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		IdleTimeout:       time.Minute,
	}

	shutdownErr := make(chan error)

	go func() {
		// trap sigterm or interupt and gracefully shutdown the server
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)

		// Block until a signal is received.
		sig := <-quit
		app.logger.PrintInfo("shuttding down", map[string]string{
			"signal": sig.String(),
		})
		// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		shutdownErr <- srv.Shutdown(ctx)
	}()

	// logger.Printf("starting %s server on :%d ", cfg.env, cfg.port)
	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownErr
	if err != nil {
		return err
	}

	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
