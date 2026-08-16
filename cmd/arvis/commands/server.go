package commands

import (
	"fmt"
	"log/slog"
	"os"
	"net/http"
	"github.com/spf13/cobra"
	"github.com/kevinkiplangat432/arvis/internal/api"
	"github.com/kevinkiplangat432/arvis/internal/store"
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
	
	// Call internal store package instead of pgxpool directly
    db, err := store.Connect(cfg.DatabaseURL)
	if err != nil {
    // If internal/store/db.go fails, we catch it here and crash the app safely
    return fmt.Errorf("failed to initialize store: %w", err)
	}
	defer db.Close() // Keep the lifecycle management here in main/server start

	logger.Info("database connected")


	srvRouter := api.NewRouter()

	logger.Info("api listening", "addr", cfg.APIAddr)
	if err := http.ListenAndServe(cfg.APIAddr, srvRouter); err != nil {
		logger.Error("api server error", "error", err)
		return fmt.Errorf("api server error: %w", err)
	}

	logger.Info("shutdown complete")
	return nil
}