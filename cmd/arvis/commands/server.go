package commands

import (
	"fmt"
	"log/slog"
	"os"
	"net/http"
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevinkiplangat432/arvis/internal/api"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Short: "Start the ARVIS proxy and API servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(banner)
		return startServer()
	},
}

func startServer() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := cfg.RequireAPIKey(); err != nil {
		return err
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	logger.Info("database connected")

	srvRouter := api.NewRouter()

	if err := http.ListenAndServe(":8080", srvRouter); err != nil {
		logger.Error("api server error", "error", err)
		os.Exit(1)
	}

	logger.Info("shutdown complete")
	return nil
}