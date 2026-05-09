//go:build integration

package repo

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

// helpers to seed data via raw SQL so repo tests can exercise the repo methods themselves
func seedUserRow(t *testing.T, ctx context.Context, id int64, locale, currency, timezone string) {
	t.Helper()
	_, err := testPool.Exec(ctx, `INSERT INTO app_user (id, locale, currency, timezone) VALUES ($1, $2, $3, $4)`, id, locale, currency, timezone)
	if err != nil {
		t.Fatalf("seed user %d: %v", id, err)
	}
}

func seedTxnRow(t *testing.T, ctx context.Context, dt string, catId int, desc string, userId int64, amount int64, currency string) {
	t.Helper()
	_, err := testPool.Exec(ctx, `INSERT INTO transaction (datetime, category_id, description, user_id, amount, currency) VALUES ($1, $2, $3, $4, $5, $6)`,
		dt, catId, desc, userId, amount, currency)
	if err != nil {
		t.Fatalf("seed txn: %v", err)
	}
}
