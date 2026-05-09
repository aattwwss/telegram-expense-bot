package domain

import (
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/entity"
)

func TestTransactionsGetFormattedHTMLMsg(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Singapore")
	dt, _ := time.ParseInLocation("2006-01-02 15:04", "2023-01-15 12:30", loc)

	trxs := Transactions{
		{
			Id:           1,
			Datetime:     dt,
			CategoryName: "Food",
			Description:  "Chicken Rice",
			Amount:       money.New(550, "SGD"),
		},
		{
			Id:           2,
			Datetime:     dt,
			CategoryName: "Transport",
			Description:  "MRT",
			Amount:       money.New(120, "SGD"),
		},
	}

	html := trxs.GetFormattedHTMLMsg(time.January, 2023, loc, 5, 0, 10)

	if len(html) == 0 {
		t.Error("expected non-empty HTML message")
	}
	if !contains(html, "January") || !contains(html, "2023") {
		t.Error("expected header with January 2023")
	}
	if !contains(html, "Chicken Rice") || !contains(html, "Food") {
		t.Error("expected first transaction details")
	}
	if !contains(html, "MRT") || !contains(html, "Transport") {
		t.Error("expected second transaction details")
	}
	if !contains(html, "[1/1]") {
		t.Error("expected page indicator [1/1]")
	}
}

func TestBreakdownsGetFormattedHTMLMsg(t *testing.T) {
	bds := Breakdowns{
		{CategoryName: "Food", Amount: money.New(5000, "SGD"), Percent: 50.0},
		{CategoryName: "Transport", Amount: money.New(3000, "SGD"), Percent: 30.0},
		{CategoryName: "Shopping", Amount: money.New(2000, "SGD"), Percent: 20.0},
	}

	html := bds.GetFormattedHTMLMsg()
	if len(html) == 0 {
		t.Error("expected non-empty HTML message")
	}
	if !contains(html, "Food") || !contains(html, "50.0") {
		t.Error("expected Food with 50.0%")
	}
	if !contains(html, "Transport") || !contains(html, "30.0") {
		t.Error("expected Transport with 30.0%")
	}
	// values < 10% should have leading space
	if !contains(html, "20.0") {
		t.Error("expected Shopping with 20.0%")
	}
}

func TestEmptyBreakdowns(t *testing.T) {
	bds := Breakdowns{}
	html := bds.GetFormattedHTMLMsg()
	if html != "" {
		t.Errorf("expected empty HTML for empty breakdowns, got %q", html)
	}
}

func TestTransactionFromEntity(t *testing.T) {
	e := entity.Transaction{
		Id:           1,
		Datetime:     time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
		CategoryId:   3,
		CategoryName: "Food",
		Description:  "Chicken Rice",
		UserId:       100,
		Amount:       550,
		Currency:     "SGD",
	}

	got := TransactionFromEntity(e)
	if got.Id != 1 {
		t.Errorf("Id = %d, want 1", got.Id)
	}
	if got.CategoryName != "Food" {
		t.Errorf("CategoryName = %s, want Food", got.CategoryName)
	}
	if got.Amount.Amount() != 550 {
		t.Errorf("Amount = %d, want 550", got.Amount.Amount())
	}
	if got.Amount.Currency().Code != "SGD" {
		t.Errorf("Currency = %s, want SGD", got.Amount.Currency().Code)
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
