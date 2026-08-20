package commands

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/kevinkiplangat432/arvis/internal/api"
	"github.com/kevinkiplangat432/arvis/internal/proxy"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

var serverCmd = &cobra.Command{
	Use:   "server",
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

	if err := cfg.ResolveProviders(); err != nil {
		return err
	}

	db, err := store.Connect(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}
	defer db.Close()

	logger.Info("database connected")
	logger.Info("providers loaded", "count", len(cfg.Providers))

	apiRouter := api.NewRouter(db, logger)
	proxyHandler := proxy.New(cfg, db, logger)

	// errgroup so either server dying brings the whole process down —
	// a proxy silently running without its API, or vice versa, is a
	// worse failure mode than both going down together and getting
	// restarted by whatever's supervising the process (systemd, Docker).
	var g errgroup.Group

	g.Go(func() error {
		logger.Info("api listening", "addr", cfg.APIAddr)
		if err := http.ListenAndServe(cfg.APIAddr, apiRouter); err != nil {
			return fmt.Errorf("api server error: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		logger.Info("proxy listening", "addr", cfg.ProxyAddr)
		if err := http.ListenAndServe(cfg.ProxyAddr, proxyHandler); err != nil {
			return fmt.Errorf("proxy server error: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Error("server error", "error", err)
		return err
	}

	logger.Info("shutdown complete")
	return nil
}