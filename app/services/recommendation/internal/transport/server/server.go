package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/AndroSaal/RecommendationsForUsers/app/services/recommendation/pkg/config"
	"github.com/pkg/errors"
)

const megabyte = 1 << 20

type Server struct {
	Server *http.Server
	Logger *slog.Logger
}

// создание нового сервера
func NewServer(cfg config.ServerConfig, handler http.Handler, logger *slog.Logger) (*Server, error) {
	server := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        handler,
		MaxHeaderBytes: megabyte,
		ReadTimeout:    cfg.Timeout,
		WriteTimeout:   cfg.Timeout,
	}

	return &Server{
		Server: server,
		Logger: logger,
	}, nil
}

// запуск сервера
func (s *Server) Run() error {
	fi := "transport.Server.Run"

	s.Logger.Info(fi + ":" + "starting server...")
	s.Logger.Info(fi + ":" + "server started on port " + s.Server.Addr)

	if err := s.Server.ListenAndServe(); err != nil {
		return errors.New(fi + ":" + "cannot run server: " + err.Error())

	}

	return nil

}

// остановка сервера
func (s *Server) Stop(ctx context.Context) {
	fi := "transport.Server.Stop"

	if err := s.Server.Shutdown(ctx); err != nil {
		s.Logger.Error(fi + ":" + "cannot stop server: " + err.Error())
	}
}
