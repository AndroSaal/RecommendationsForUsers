package server

import (
	"context"
	"net/http"
)

type Server struct {
	Server *http.Server
	logger *slog.Logger
}

// создание нового сервера
func NewServer(cfg config.SrvConfig, handler http.Handler, logger *slog.Logger) (*Server, error) {
	server := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    cfg.Timeout,
		WriteTimeout:   cfg.Timeout,
	}

	return &Server{
		Server: server,
		logger: logger,
	}, nil
}

// запуск сервера
func (s *Server) MustRun() {
	fi := "transport.Server.MustRun"

	s.logger.Info(fi + ":" + "starting server...")
	s.logger.Info(fi + ":" + "server started on port " + s.Server.Addr)

	if err := s.Server.ListenAndServe(); err != nil {
		panic(fi + ":" + "cannot run server: " + err.Error())
	}

}

// остановка сервера
func (s *Server) Stop(ctx context.Context) {
	fi := "transport.Server.Stop"

	if err := s.Server.Shutdown(ctx); err != nil {
		s.logger.Error(fi + ":" + "cannot stop server: " + err.Error())
	}
}