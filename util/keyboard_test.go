package util

import (
	"encoding/json"
	"testing"

	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/enum"
)

func TestNewInlineKeyboard(t *testing.T) {
	configs := []InlineKeyboardConfig{
		{label: "A", data: "data-a"},
		{label: "B", data: "data-b"},
		{label: "C", data: "data-c"},
	}

	// 3 items, 2 cols, cancellable → 2 data rows + 1 cancel = 3
	kb := NewInlineKeyboard(configs, 1, 2, true)
	if len(kb) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(kb))
	}
	if len(kb[0]) != 2 || len(kb[1]) != 1 {
		t.Errorf("expected row 0 len 2, row 1 len 1, got %d and %d", len(kb[0]), len(kb[1]))
	}
	if kb[0][0].Text != "A" || kb[0][1].Text != "B" || kb[1][0].Text != "C" {
		t.Errorf("unexpected button labels")
	}

	// non-cancellable → 2 data rows = 2
	kb = NewInlineKeyboard(configs, 1, 2, false)
	if len(kb) != 2 {
		t.Errorf("expected 2 rows without cancel, got %d", len(kb))
	}

	// empty configs with cancel → cancel row only = 1
	kb = NewInlineKeyboard(nil, 1, 2, true)
	if len(kb) != 1 {
		t.Fatalf("expected 1 row, got %d", len(kb))
	}
	if kb[0][0].Text != "Cancel" {
		t.Errorf("expected Cancel button, got %q", kb[0][0].Text)
	}
}

func TestNewPaginationKeyboard(t *testing.T) {
	// middle of list: should have prev, next, and cancel
	kb, err := NewPaginationKeyboard(100, 10, 10, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(kb) != 2 { // 1 data row + cancel = 2
		t.Fatalf("expected 2 rows, got %d", len(kb))
	}
	if kb[0][0].Text != "<" || kb[0][1].Text != ">" {
		t.Errorf("expected '<' and '>' buttons, got %q, %q", kb[0][0].Text, kb[0][1].Text)
	}
	if kb[1][0].Text != "Cancel" {
		t.Errorf("expected Cancel button in last row")
	}

	// first page: only next button + cancel
	kb, err = NewPaginationKeyboard(100, 0, 10, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(kb[0]) != 1 || kb[0][0].Text != ">" {
		t.Errorf("expected only '>' button, got row len %d", len(kb[0]))
	}

	// only one page: totalCount 8 < limit 10, no next, no prev
	// len(configs) == 0 so showCancelButton = false → empty result
	kb, err = NewPaginationKeyboard(8, 0, 10, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(kb) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(kb))
	}
}

func TestNewUndoConfirmationKeyboard(t *testing.T) {
	kb, err := NewUndoConfirmationKeyboard(42, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(kb) != 2 { // Yes row + Cancel row = 2
		t.Fatalf("expected 2 rows, got %d", len(kb))
	}
	if kb[0][0].Text != "Yes" {
		t.Errorf("expected 'Yes' button, got %q", kb[0][0].Text)
	}

	// Parse callback data to verify transaction id
	var undo domain.UndoCallback
	err = json.Unmarshal([]byte(*kb[0][0].CallbackData), &undo)
	if err != nil {
		t.Fatalf("unexpected error unmarshalling callback data: %v", err)
	}
	if undo.TransactionId != 42 {
		t.Errorf("expected TransactionId 42, got %d", undo.TransactionId)
	}
	if undo.Type != enum.Undo {
		t.Errorf("expected Type Undo, got %v", undo.Type)
	}

	// Last row should have Cancel
	lastRow := kb[len(kb)-1]
	if lastRow[0].Text != "Cancel" {
		t.Errorf("expected 'Cancel' button, got %q", lastRow[0].Text)
	}
}
