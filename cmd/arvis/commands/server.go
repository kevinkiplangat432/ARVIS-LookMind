package commands

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

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevinkiplangat432/arvis/internal/api"
	"github.com/kevinkiplangat432/arvis/internal/config"
	"github.com/kevinkiplangat432/arvis/internal/proxy"
	"github.com/spf13/cobra"
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

	cfg := config.Load()

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

	m, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to init migrations: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.Info("migrations applied")

	targetURL, err := url.Parse(cfg.TargetURL)
	if err != nil {
		return fmt.Errorf("invalid target URL %q: %w", cfg.TargetURL, err)
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
	return nil
}




// code  review  ((purposes))

// initialization and cofiguration
// 
// lines 30 - 34 (slog.NewJSONHandler) ..structured logs are needed for soc 2 sudit trails. it outputs logs in machine-readable JSON

// line 36 (config.Load()): centralizes configuration management 
// lines 38-39 (context.WithTimeouts) prevents the application from hanging indefinitely during startup if infrastructure is down

// database infra
//
//line 41 - 45 (pgxpool.New) pgxpool handles connection pooling ..this is critical for 1000 + RPS without exhausting database file descriptors
// line 47 - 50 (db.Ping) ensure the data base is fully reachable  beforelaunching network listeners
//lines 52 - 58 (migrate.New / m.up()) runs schema migrationsautomatically.while good for dev, executing mugration automatically on server startup introduces serious risk at production scale 


//proxy and httpps servers 
// lines 60 - 65 (url.Parse / newSingleHostReverseProxy) Pre-parse the backend destination URLand initializes Go's standard reverse proxy
//lines 68 - 80(proxyserver & apiserver ) instantiates two separate Http servers with explicit read, write and idle timeouts. These timeouts prevent Slowloris Dos attacks
//lines 82 -96 ( go func() ListenAndServe): spawns concurrent background goroutines so both servers can run simultaneously without blocking each other



//gracefull shutdown
//lines 98 - 100 (signal.Notify) traps the system killsignals(SIGIT,SIGTERM) to prevent hard crashes
// lines 104-108 (shutdown(shutctx)) allows active HTTP request up to 15 seconds to finsih processing before the process exits cleanly




// code review ii (Flaws, Optimizations & SOC 2 Gap Analysis)

// data leakage & information Disclosure 
// The Issue: fmt.Errorf("invalid target URL %q: %w", cfg.TargetURL, err) and database connection errors are bubbled up or printed. If cfg.DatabaseURL contains a password embedded in the connection string, any failure to initialize migrations or pools could print raw credentials to os.Stderr or log systems.The Fix: Sanitize URLs before logging them or wrapping them in errors. Strip passwords out of connection strings using url.Parse and replacing the user info password field before logging.

// soc 2 type ii Compliance Gaps 
// Automated Migrations on Startup: Running m.Up() inside your application binary breaks the principle of Least Privilege. The runtime database user for a proxy only needs SELECT/INSERT/UPDATE permissions. To run migrations, it needs DDL permissions (CREATE TABLE, ALTER). If the application is compromised, an attacker gains full control over the schema.Remediation: Move migrations to a separate CI/CD deployment step or a dedicated migration job.Unstructured/Uncontrolled Panics: os.Exit(1) inside the goroutines abruptly kills the process. This skips your graceful shutdown code (defer db.Close() and server .Shutdown()), which can corrupt state or drop active sessions—violating SOC 2 Availability principles.



// scaling to 1000 RPS 
// Default Connection Pool Limits: pgxpool.New uses standard default connection limits. At 1,000 RPS, concurrent goroutines may quickly exhaust the pool, causing bottleneck latencies or strict context deadlines.Global Logging Contention: slog.NewJSONHandler(os.Stdout, ...) writes directly to standard output synchronous-style. At high throughput, thread contention on os.Stdout becomes a heavy performance bottleneck.Reverse Proxy Fine-Tuning: The default httputil.ReverseProxy uses http.DefaultTransport. This transport maintains shared global idle connection pools that are not tuned for 1,000+ continuous RPS, which will result in socket exhaustion (TIME_WAIT tokens).