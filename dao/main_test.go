//go:build integration

package dao

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/aattwwss/telegram-expense-bot/internal/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	pool, cleanup, err := testutil.StartPostgres(ctx, "../scripts/sql")
	if err != nil {
		log.Fatalf("start postgres: %v", err)
	}
	testPool = pool
	code := m.Run()
	cleanup()
	testPool.Close()
	os.Exit(code)
}

func clearTables(t *testing.T, ctx context.Context) {
	t.Helper()
	tables := []string{"transaction", "message_context", "app_user"}
	for _, table := range tables {
		if _, err := testPool.Exec(ctx, "DELETE FROM "+table); err != nil {
			t.Fatalf("clear %s: %v", table, err)
		}
	}
}
