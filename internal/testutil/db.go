package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func init() {
	// Auto-detect Podman socket if DOCKER_HOST is not already set.
	// Docker is auto-detected by testcontainers; Podman requires DOCKER_HOST.
	if os.Getenv("DOCKER_HOST") == "" {
		candidates := []string{
			"/run/user/" + fmt.Sprint(os.Getuid()) + "/podman/podman.sock",
			"/var/run/podman/podman.sock",
		}
		for _, p := range candidates {
			if _, err := os.Stat(p); err == nil {
				os.Setenv("DOCKER_HOST", "unix://"+p)
				return
			}
		}
	}
}

const (
	defaultImage = "postgres:16-alpine"
	defaultUser  = "test"
	defaultPass  = "test"
	defaultDB    = "testdb"
)

// StartPostgres starts a postgres container, runs init.sql, and returns a connection pool.
// The caller is responsible for calling the returned cleanup function.
func StartPostgres(ctx context.Context, initScriptPath string) (*pgxpool.Pool, func(), error) {
	absPath, err := filepath.Abs(initScriptPath)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve init script path: %w", err)
	}
	if _, err := os.Stat(absPath); err != nil {
		return nil, nil, fmt.Errorf("init script not found at %s: %w", absPath, err)
	}

	ctr, err := postgres.Run(ctx, defaultImage,
		postgres.WithUsername(defaultUser),
		postgres.WithPassword(defaultPass),
		postgres.WithDatabase(defaultDB),
		postgres.WithInitScripts(absPath),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("start postgres container: %w", err)
	}

	cleanup := func() {
		if err := ctr.Terminate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to terminate postgres container: %v\n", err)
		}
	}

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable", "search_path=public")
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("get connection string: %w", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		cleanup()
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, cleanup, nil
}
