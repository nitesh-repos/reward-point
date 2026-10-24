package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"reward-points/internal/database"
	"reward-points/internal/server"
	"syscall"
	"time"
)

func main() {
	cfg := server.Config{
		Host: "0.0.0.0",
		Port: 8080,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	serv := server.New(cfg, logger, database.NewLocalStorage())

	signalChan := make(chan os.Signal, 1)
	defer signal.Stop(signalChan)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		if err := serv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	case <-signalChan:
		logger.Info("interrupt signal received. shutting down server")
		ctx, canFunc := context.WithTimeout(context.Background(), time.Second*10)
		defer canFunc()
		if err := serv.Stop(ctx); err != nil {
			logger.Error("error shutting down server", "error", err)
			os.Exit(1)
		}
	}

}
