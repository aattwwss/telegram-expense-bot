//go:build integration

package repo

import (
	"context"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/entity"
)

func newTestTransactionRepo() TransactionRepo {
	return NewTransactionRepo(dao.NewTransactionDao(testPool))
}

func TestTransactionRepo_Add(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUserRow(t, ctx, 100, "en", "SGD", "Asia/Singapore")

	repo := newTestTransactionRepo()

	dt := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	err := repo.Add(ctx, domain.Transaction{
		Datetime:     dt,
		CategoryId:   1,
		CategoryName: "Bills",
		Description:  "Electric bill",
		UserId:       100,
		Amount:       money.New(4500, "SGD"),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	// Verify via DAO directly
	trx, err := dao.NewTransactionDao(testPool).FindLatestByUserId(ctx, 100)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if trx.Amount != 4500 {
		t.Errorf("Amount = %d, want 4500", trx.Amount)
	}
	if trx.Currency != "SGD" {
		t.Errorf("Currency = %s, want SGD", trx.Currency)
	}
}

func TestTransactionRepo_GetById(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUserRow(t, ctx, 100, "en", "SGD", "Asia/Singapore")
	seedTxnRow(t, ctx, "2024-03-15T10:30:00Z", 2, "MRT ride", 100, 120, "SGD")

	repo := newTestTransactionRepo()

	// Get the ID via DAO
	latest, _ := dao.NewTransactionDao(testPool).FindLatestByUserId(ctx, 100)

	got, err := repo.GetById(ctx, latest.Id, 100)
	if err != nil {
		t.Fatalf("GetById: %v", err)
	}
	if got.Amount.Amount() != 120 {
		t.Errorf("Amount = %d, want 120", got.Amount.Amount())
	}
	if got.CategoryName != "Education" {
		t.Errorf("CategoryName = %s, want Education", got.CategoryName)
	}
}

func TestTransactionRepo_FindLastestByUserId(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUserRow(t, ctx, 100, "en", "SGD", "Asia/Singapore")
	seedTxnRow(t, ctx, "2024-01-01T10:00:00Z", 1, "old", 100, 100, "SGD")
	seedTxnRow(t, ctx, "2024-06-15T15:00:00Z", 1, "new", 100, 500, "SGD")

	repo := newTestTransactionRepo()

	latest, err := repo.FindLastestByUserId(ctx, 100)
	if err != nil {
		t.Fatalf("FindLastestByUserId: %v", err)
	}
	if latest == nil {
		t.Fatal("expected transaction, got nil")
	}
	if latest.Amount.Amount() != 500 {
		t.Errorf("Amount = %d, want 500", latest.Amount.Amount())
	}
}

func TestTransactionRepo_FindLastestByUserId_Empty(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUserRow(t, ctx, 100, "en", "SGD", "Asia/Singapore")

	repo := newTestTransactionRepo()

	latest, err := repo.FindLastestByUserId(ctx, 100)
	if err != nil {
		t.Fatalf("FindLastestByUserId: %v", err)
	}
	if latest != nil {
		t.Errorf("expected nil, got %+v", latest)
	}
}

func TestTransactionRepo_DeleteById(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUserRow(t, ctx, 100, "en", "SGD", "Asia/Singapore")
	seedTxnRow(t, ctx, "2024-04-10T08:00:00Z", 1, "to delete", 100, 300, "SGD")

	repo := newTestTransactionRepo()
	trxDAO := dao.NewTransactionDao(testPool)

	latest, _ := trxDAO.FindLatestByUserId(ctx, 100)

	if err := repo.DeleteById(ctx, latest.Id, 100); err != nil {
		t.Fatalf("DeleteById: %v", err)
	}

	_, err := trxDAO.GetById(ctx, latest.Id, 100)
	if err == nil {
		t.Fatal("expected not found after delete")
	}
}

func TestTransactionRepo_GetTransactionBreakdownByCategory(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)

	loc, _ := time.LoadLocation("Asia/Singapore")
	user := domain.User{
		Id:       100,
		Locale:   "en",
		Currency: money.GetCurrency("SGD"),
		Location: loc,
	}
	seedUserRow(t, ctx, 100, "en", "SGD", "Asia/Singapore")

	// June 2024: 2 Food (cat 4) = 500+300=800, 1 Transport (cat 13) = 200
	seedTxnRow(t, ctx, "2024-06-01T10:00:00+08:00", 4, "lunch", 100, 500, "SGD")
	seedTxnRow(t, ctx, "2024-06-10T12:00:00+08:00", 4, "dinner", 100, 300, "SGD")
	seedTxnRow(t, ctx, "2024-06-20T14:00:00+08:00", 13, "bus", 100, 200, "SGD")

	repo := newTestTransactionRepo()

	breakdowns, total, err := repo.GetTransactionBreakdownByCategory(ctx, time.June, 2024, user)
	if err != nil {
		t.Fatalf("GetTransactionBreakdownByCategory: %v", err)
	}
	if len(breakdowns) != 2 {
		t.Fatalf("len = %d, want 2", len(breakdowns))
	}
	if total.Amount() != 1000 {
		t.Errorf("total = %d, want 1000", total.Amount())
	}

	// Food = 80.0%, Transport = 20.0%
	if breakdowns[0].CategoryName != "Food" {
		t.Errorf("first = %s, want Food", breakdowns[0].CategoryName)
	}
	if breakdowns[0].Percent != 80.0 {
		t.Errorf("Food percent = %f, want 80.0", breakdowns[0].Percent)
	}
	if breakdowns[1].CategoryName != "Transport" {
		t.Errorf("second = %s, want Transport", breakdowns[1].CategoryName)
	}
	if breakdowns[1].Percent != 20.0 {
		t.Errorf("Transport percent = %f, want 20.0", breakdowns[1].Percent)
	}
}

func TestTransactionRepo_ListByMonthAndYear(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUserRow(t, ctx, 100, "en", "SGD", "Asia/Singapore")

	// June 2024: 3 transactions
	seedTxnRow(t, ctx, "2024-06-01T10:00:00+08:00", 1, "txn1", 100, 100, "SGD")
	seedTxnRow(t, ctx, "2024-06-10T12:00:00+08:00", 2, "txn2", 100, 200, "SGD")
	seedTxnRow(t, ctx, "2024-06-20T14:00:00+08:00", 3, "txn3", 100, 300, "SGD")
	// Outside range
	seedTxnRow(t, ctx, "2024-07-01T00:00:00+08:00", 1, "july", 100, 400, "SGD")

	loc, _ := time.LoadLocation("Asia/Singapore")
	repo := newTestTransactionRepo()

	q := entity.TransactionListQuery{
		Month:    time.June,
		Year:     2024,
		Offset:   0,
		Limit:    2,
		Asc:      false,
		UserId:   100,
		Location: loc,
	}

	trxs, total, err := repo.ListByMonthAndYear(ctx, q)
	if err != nil {
		t.Fatalf("ListByMonthAndYear: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(trxs) != 2 {
		t.Fatalf("len = %d, want 2", len(trxs))
	}
	// DESC order, newest (300) first
	if trxs[0].Amount.Amount() != 300 {
		t.Errorf("first amount = %d, want 300", trxs[0].Amount.Amount())
	}
}

func TestTransactionRepo_ListByMonthAndYear_Empty(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUserRow(t, ctx, 100, "en", "SGD", "Asia/Singapore")

	loc, _ := time.LoadLocation("Asia/Singapore")
	repo := newTestTransactionRepo()

	q := entity.TransactionListQuery{
		Month:    time.January,
		Year:     2024,
		Offset:   0,
		Limit:    10,
		Asc:      false,
		UserId:   100,
		Location: loc,
	}

	trxs, total, err := repo.ListByMonthAndYear(ctx, q)
	if err != nil {
		t.Fatalf("ListByMonthAndYear: %v", err)
	}
	if total != 0 {
		t.Errorf("total = %d, want 0", total)
	}
	if len(trxs) != 0 {
		t.Errorf("len = %d, want 0", len(trxs))
	}
}
