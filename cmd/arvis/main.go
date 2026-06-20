package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kevinkiplangat432/arvis/internal/api"
	"github.com/kevinkiplangat432/arvis/internal/config"
	"github.com/kevinkiplangat432/arvis/internal/proxy"
)

const banner = `
 █████╗  ██████╗  ██╗   ██╗ ██╗ ███████╗
██╔══██╗ ██╔══██╗ ██║   ██║ ██║ ██╔════╝
███████║ ██████╔╝ ██║   ██║ ██║ ███████╗
██╔══██║ ██╔══██╗ ╚██╗ ██╔╝ ██║ ╚════██║
██║  ██║ ██║  ██║  ╚████╔╝  ██║ ███████║
╚═╝  ╚═╝ ╚═╝  ╚═╝   ╚═══╝   ╚═╝ ╚══════╝

Automated Runtime Visibility & Intelligence System
Version 0.2.0
`

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	fmt.Print(banner)

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		logger.Error("database ping failed", "error", err)
		os.Exit(1)
	}
	logger.Info("database connected")

	// Proxy server on :8080
	targetURL, err := url.Parse(cfg.TargetURL)
	if err != nil {
		logger.Error("invalid target URL", "error", err, "target", cfg.TargetURL)
		os.Exit(1)
	}
	rp := httputil.NewSingleHostReverseProxy(targetURL)
	proxyHandler := proxy.New(rp, cfg, db, logger)

	proxyServer := &http.Server{
		Addr:         cfg.ProxyAddr,
		Handler:      proxyHandler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// API server on :8081
	apiServer := &http.Server{
		Addr:         cfg.APIAddr,
		Handler:      api.NewRouter(db, logger),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("proxy listening", "addr", cfg.ProxyAddr)
		if err := proxyServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("proxy server error", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		logger.Info("api listening", "addr", cfg.APIAddr)
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("api server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down")

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutCancel()

	_ = proxyServer.Shutdown(shutCtx)
	_ = apiServer.Shutdown(shutCtx)

	logger.Info("shutdown complete")
}