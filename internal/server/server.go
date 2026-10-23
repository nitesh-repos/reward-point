package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type Config struct {
	Host string
	Port int
}

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(cfg Config, logger *slog.Logger) *Server {
	serv := &Server{
		logger: logger,
	}
	mux := http.NewServeMux()

	serv.registerRoutes(mux)

	serv.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: mux,
	}

	return serv
}

func (serv *Server) Start() error {
	serv.logger.Info("Starting server", "addr", serv.httpServer.Addr)
	return serv.httpServer.ListenAndServe()
}

func (serv *Server) Stop(ctx context.Context) error {
	return serv.httpServer.Shutdown(ctx)
}
