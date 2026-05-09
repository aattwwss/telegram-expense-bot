package domain

import (
	"testing"

	"github.com/aattwwss/telegram-expense-bot/entity"
)

func TestUserFromEntity(t *testing.T) {
	e := entity.User{
		Id:       123,
		Locale:   "en",
		Currency: "SGD",
		Timezone: "Asia/Singapore",
	}

	got, err := UserFromEntity(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Id != 123 {
		t.Errorf("Id = %d, want 123", got.Id)
	}
	if got.Locale != "en" {
		t.Errorf("Locale = %s, want en", got.Locale)
	}
	if got.Currency.Code != "SGD" {
		t.Errorf("Currency = %s, want SGD", got.Currency.Code)
	}
	if got.Location.String() != "Asia/Singapore" {
		t.Errorf("Location = %s, want Asia/Singapore", got.Location.String())
	}
}

func TestUserFromEntity_BadTimezone(t *testing.T) {
	e := entity.User{
		Id:       1,
		Currency: "SGD",
		Timezone: "Mars/Nonexistent",
	}

	_, err := UserFromEntity(e)
	if err == nil {
		t.Error("expected error for invalid timezone")
	}
}
