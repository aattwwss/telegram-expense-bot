//go:build integration

package dao

import (
	"context"
	"testing"
	"time"

	"github.com/aattwwss/telegram-expense-bot/entity"
)

func seedUser(t *testing.T, ctx context.Context, id int64) {
	t.Helper()
	userDAO := NewUserDao(testPool)
	if err := userDAO.Insert(ctx, entity.User{
		Id: id, Locale: "en", Currency: "SGD", Timezone: "Asia/Singapore",
	}); err != nil {
		t.Fatalf("seed user %d: %v", id, err)
	}
}

func insertTxn(t *testing.T, ctx context.Context, dao TransactionDAO, dt time.Time, catId int, desc string, userId int64, amount int64, currency string) int {
	t.Helper()
	err := dao.Insert(ctx, entity.Transaction{
		Datetime:   dt,
		CategoryId: catId,
		Description: desc,
		UserId:     userId,
		Amount:     amount,
		Currency:   currency,
	})
	if err != nil {
		t.Fatalf("insert txn: %v", err)
	}
	latest, err := dao.FindLatestByUserId(ctx, userId)
	if err != nil {
		t.Fatalf("get inserted id: %v", err)
	}
	return latest.Id
}

func TestTransactionDAO_InsertAndGetById(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUser(t, ctx, 100)

	dao := NewTransactionDao(testPool)

	dt := time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)
	id := insertTxn(t, ctx, dao, dt, 1, "Chicken Rice", 100, 550, "SGD")

	got, err := dao.GetById(ctx, id, 100)
	if err != nil {
		t.Fatalf("GetById: %v", err)
	}
	if got.Amount != 550 {
		t.Errorf("Amount = %d, want 550", got.Amount)
	}
	if got.Currency != "SGD" {
		t.Errorf("Currency = %s, want SGD", got.Currency)
	}
	if got.CategoryName != "Bills" {
		t.Errorf("CategoryName = %s, want Bills (from JOIN)", got.CategoryName)
	}
	if got.Description != "Chicken Rice" {
		t.Errorf("Description = %s, want Chicken Rice", got.Description)
	}
}

func TestTransactionDAO_GetById_NotFound(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUser(t, ctx, 100)

	dao := NewTransactionDao(testPool)
	_, err := dao.GetById(ctx, 99999, 100)
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestTransactionDAO_FindLatestByUserId(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUser(t, ctx, 100)

	dao := NewTransactionDao(testPool)

	older := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	newer := time.Date(2024, 6, 15, 15, 0, 0, 0, time.UTC)

	insertTxn(t, ctx, dao, older, 2, "old txn", 100, 100, "SGD")
	insertTxn(t, ctx, dao, newer, 3, "new txn", 100, 200, "SGD")

	latest, err := dao.FindLatestByUserId(ctx, 100)
	if err != nil {
		t.Fatalf("FindLatestByUserId: %v", err)
	}
	if latest == nil {
		t.Fatal("expected transaction, got nil")
	}
	if latest.Amount != 200 {
		t.Errorf("Amount = %d, want 200 (newer)", latest.Amount)
	}
}

func TestTransactionDAO_FindLatestByUserId_Empty(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUser(t, ctx, 100)

	dao := NewTransactionDao(testPool)
	latest, err := dao.FindLatestByUserId(ctx, 100)
	if err != nil {
		t.Fatalf("FindLatestByUserId: %v", err)
	}
	if latest != nil {
		t.Errorf("expected nil, got %+v", latest)
	}
}

func TestTransactionDAO_DeleteById(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUser(t, ctx, 100)

	dao := NewTransactionDao(testPool)
	dt := time.Date(2024, 3, 10, 8, 0, 0, 0, time.UTC)
	id := insertTxn(t, ctx, dao, dt, 1, "to delete", 100, 300, "SGD")

	if err := dao.DeleteById(ctx, id, 100); err != nil {
		t.Fatalf("DeleteById: %v", err)
	}

	_, err := dao.GetById(ctx, id, 100)
	if err == nil {
		t.Fatal("expected not found after delete")
	}
}

func TestTransactionDAO_CountAndListByMonthAndYear(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUser(t, ctx, 100)
	dao := NewTransactionDao(testPool)

	// Insert 3 transactions in June 2024
	insertTxn(t, ctx, dao, time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC), 1, "txn 1", 100, 100, "SGD")
	insertTxn(t, ctx, dao, time.Date(2024, 6, 10, 12, 0, 0, 0, time.UTC), 2, "txn 2", 100, 200, "SGD")
	insertTxn(t, ctx, dao, time.Date(2024, 6, 20, 14, 0, 0, 0, time.UTC), 3, "txn 3", 100, 300, "SGD")

	// Also insert one outside the range
	insertTxn(t, ctx, dao, time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), 1, "july txn", 100, 400, "SGD")

	dateFrom := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)

	count, err := dao.CountListByMonthAndYear(ctx, dateFrom, dateTo, 100)
	if err != nil {
		t.Fatalf("CountListByMonthAndYear: %v", err)
	}
	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}

	// Test pagination: offset 0, limit 2
	results, err := dao.ListByMonthAndYear(ctx, dateFrom, dateTo, 0, 2, false, 100)
	if err != nil {
		t.Fatalf("ListByMonthAndYear: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("len = %d, want 2", len(results))
	}
	// DESC order, so newest first
	if results[0].Amount != 300 {
		t.Errorf("first Amount (DESC) = %d, want 300", results[0].Amount)
	}

	// Test ascending order
	ascResults, err := dao.ListByMonthAndYear(ctx, dateFrom, dateTo, 0, 10, true, 100)
	if err != nil {
		t.Fatalf("ListByMonthAndYear asc: %v", err)
	}
	if len(ascResults) != 3 {
		t.Fatalf("len asc = %d, want 3", len(ascResults))
	}
	if ascResults[0].Amount != 100 {
		t.Errorf("first Amount (ASC) = %d, want 100", ascResults[0].Amount)
	}
}

func TestTransactionDAO_GetBreakdownByCategory(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)
	seedUser(t, ctx, 100)
	dao := NewTransactionDao(testPool)

	// Insert 2 Food (cat 4) and 1 Transport (cat 13) in June 2024
	insertTxn(t, ctx, dao, time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC), 4, "lunch", 100, 500, "SGD")      // cat 4 = Food
	insertTxn(t, ctx, dao, time.Date(2024, 6, 10, 12, 0, 0, 0, time.UTC), 4, "dinner", 100, 300, "SGD")    // cat 4 = Food
	insertTxn(t, ctx, dao, time.Date(2024, 6, 20, 14, 0, 0, 0, time.UTC), 13, "bus", 100, 200, "SGD")      // cat 13 = Transport

	dateFrom := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)

	breakdowns, err := dao.GetBreakdownByCategory(ctx, dateFrom, dateTo, 100)
	if err != nil {
		t.Fatalf("GetBreakdownByCategory: %v", err)
	}
	if len(breakdowns) != 2 {
		t.Fatalf("len = %d, want 2", len(breakdowns))
	}
	// Ordered by amount DESC; Food total = 800, Transport = 200
	if breakdowns[0].CategoryName != "Food" {
		t.Errorf("first category = %s, want Food", breakdowns[0].CategoryName)
	}
	if breakdowns[0].Amount != 800 {
		t.Errorf("Food amount = %d, want 800", breakdowns[0].Amount)
	}
	if breakdowns[1].CategoryName != "Transport" {
		t.Errorf("second category = %s, want Transport", breakdowns[1].CategoryName)
	}
	if breakdowns[1].Amount != 200 {
		t.Errorf("Transport amount = %d, want 200", breakdowns[1].Amount)
	}
}
