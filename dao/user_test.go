//go:build integration

package dao

import (
	"context"
	"testing"

	"github.com/aattwwss/telegram-expense-bot/entity"
)

func TestUserDAO_InsertAndFindById(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)

	dao := NewUserDao(testPool)

	err := dao.Insert(ctx, entity.User{
		Id:       12345,
		Locale:   "en",
		Currency: "SGD",
		Timezone: "Asia/Singapore",
	})
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}

	user, err := dao.FindUserById(ctx, 12345)
	if err != nil {
		t.Fatalf("FindUserById: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Id != 12345 {
		t.Errorf("Id = %d, want 12345", user.Id)
	}
	if user.Locale != "en" {
		t.Errorf("Locale = %s, want en", user.Locale)
	}
	if user.Currency != "SGD" {
		t.Errorf("Currency = %s, want SGD", user.Currency)
	}
	if user.Timezone != "Asia/Singapore" {
		t.Errorf("Timezone = %s, want Asia/Singapore", user.Timezone)
	}
}

func TestUserDAO_FindById_NotFound(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)

	dao := NewUserDao(testPool)

	user, err := dao.FindUserById(ctx, 99999)
	if err != nil {
		t.Fatalf("FindUserById: %v", err)
	}
	if user != nil {
		t.Errorf("expected nil user, got %+v", user)
	}
}
