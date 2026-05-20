package http

import (
	"net/http"
	"time"

	"github.com/ssu526/api-gateway/internal/config"
)

func NewServer(cfg *config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        handler,
		ReadTimeout:    time.Duration(cfg.ReadTimeout),
		WriteTimeout:   time.Duration(cfg.WriteTimeout),
		IdleTimeout:    time.Duration(cfg.IdleTimeout),
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}
}
