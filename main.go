package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/rohitxdev/abc-task/internal/config"
	"github.com/rohitxdev/abc-task/internal/database"
	"github.com/rohitxdev/abc-task/internal/handler"
	"github.com/rohitxdev/abc-task/internal/repo"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	db, err := database.NewSQLite(cfg.DatabaseURL)
	if err != nil {
		panic("Failed to create database: " + err.Error())
	}
	defer db.Close()

	r, err := repo.New(db)
	if err != nil {
		panic("Failed to create repo: " + err.Error())
	}

	svc := &handler.Services{
		Config: cfg,
		Repo:   r,
	}

	h, err := handler.New(svc)
	if err != nil {
		panic("Failed to create handler: " + err.Error())
	}

	ls, err := net.Listen("tcp", net.JoinHostPort(cfg.Host, cfg.Port))
	if err != nil {
		panic("Failed to listen on TCP: " + err.Error())
	}
	defer func() {
		if err = ls.Close(); err != nil {
			panic("Failed to close TCP listener: " + err.Error())
		}
	}()

	//Start HTTP server
	go func() {
		// Stdlib supports HTTP/2 by default when serving over TLS, but has to be explicitly enabled otherwise.
		handler := h2c.NewHandler(h, &http2.Server{})
		if err := http.Serve(ls, handler); err != nil && !errors.Is(err, net.ErrClosed) {
			panic("Failed to serve HTTP: " + err.Error())
		}
	}()

	slog.Debug("HTTP server started")
	slog.Info(fmt.Sprintf("Server is listening on http://%s and is ready to serve requests", ls.Addr()))

	ctx := context.Background()
	//Shut down HTTP server gracefully
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(ctx, cfg.ShutdownTimeout)
	defer cancel()

	if err := h.Shutdown(ctx); err != nil {
		panic("Failed to shutdown HTTP server: " + err.Error())
	}

	slog.Debug("HTTP server shut down gracefully")
}
