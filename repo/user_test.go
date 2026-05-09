//go:build integration

package repo

import (
	"context"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/domain"
)

func newTestUserRepo() UserRepo {
	return NewUserRepo(dao.NewUserDao(testPool))
}

func TestUserRepo_AddAndFindById(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)

	repo := newTestUserRepo()

	err := repo.Add(ctx, domain.User{
		Id:       777,
		Locale:   "ja",
		Currency: money.GetCurrency("SGD"),
		Location: timeLocation("Asia/Tokyo"),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	user, err := repo.FindUserById(ctx, 777)
	if err != nil {
		t.Fatalf("FindUserById: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Id != 777 {
		t.Errorf("Id = %d, want 777", user.Id)
	}
	if user.Locale != "ja" {
		t.Errorf("Locale = %s, want ja", user.Locale)
	}
	if user.Currency.Code != "SGD" {
		t.Errorf("Currency = %s, want SGD", user.Currency.Code)
	}
	if user.Location.String() != "Asia/Tokyo" {
		t.Errorf("Location = %s, want Asia/Tokyo", user.Location.String())
	}
}

func TestUserRepo_FindById_NotFound(t *testing.T) {
	ctx := context.Background()
	clearTables(t, ctx)

	repo := newTestUserRepo()

	user, err := repo.FindUserById(ctx, 99999)
	if err != nil {
		t.Fatalf("FindUserById: %v", err)
	}
	if user != nil {
		t.Errorf("expected nil, got %+v", user)
	}
}

func timeLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
