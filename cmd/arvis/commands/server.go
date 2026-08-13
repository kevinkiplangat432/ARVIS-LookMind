package commands

import (
	"fmt"
	"log/slog"
	"os"
	"net/http"

	"github.com/spf13/cobra"
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
	if err := cfg.RequireAPIKey(); err != nil {
		return err
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	srvRouter := api.NewRouter()

	if err := http.ListenAndServe(":8080", srvRouter); err != nil {
		logger.Error("api server error", "error", err)
		os.Exit(1)
	}

	logger.Info("shutdown complete")
	return nil
}