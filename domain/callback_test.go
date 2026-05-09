package domain

import (
	"encoding/json"
	"testing"

	"github.com/aattwwss/telegram-expense-bot/enum"
)

func TestCategoryCallback_JSON(t *testing.T) {
	original := CategoryCallback{
		Callback: Callback{
			Type:             enum.Category,
			MessageContextId: 42,
		},
		CategoryId: 7,
	}

	var decoded CategoryCallback
	roundTrip(original, &decoded, t)

	if decoded.Type != enum.Category {
		t.Errorf("Type = %v, want Category", decoded.Type)
	}
	if decoded.MessageContextId != 42 {
		t.Errorf("MessageContextId = %d, want 42", decoded.MessageContextId)
	}
	if decoded.CategoryId != 7 {
		t.Errorf("CategoryId = %d, want 7", decoded.CategoryId)
	}
}

func TestPaginationCallback_JSON(t *testing.T) {
	original := PaginationCallback{
		Callback: Callback{
			Type:             enum.Pagination,
			MessageContextId: 10,
		},
		Action: enum.Next,
		Offset: 30,
		Limit:  10,
	}

	var decoded PaginationCallback
	roundTrip(original, &decoded, t)

	if decoded.Offset != 30 {
		t.Errorf("Offset = %d, want 30", decoded.Offset)
	}
	if decoded.Limit != 10 {
		t.Errorf("Limit = %d, want 10", decoded.Limit)
	}
	if decoded.Action != enum.Next {
		t.Errorf("Action = %v, want Next", decoded.Action)
	}
}

func TestUndoCallback_JSON(t *testing.T) {
	original := UndoCallback{
		Callback: Callback{
			Type:             enum.Undo,
			MessageContextId: 5,
		},
		TransactionId: 99,
	}

	var decoded UndoCallback
	roundTrip(original, &decoded, t)

	if decoded.TransactionId != 99 {
		t.Errorf("TransactionId = %d, want 99", decoded.TransactionId)
	}
	if decoded.Type != enum.Undo {
		t.Errorf("Type = %v, want Undo", decoded.Type)
	}
}

func TestGenericCallback_JSON(t *testing.T) {
	original := GenericCallback{
		Callback: Callback{
			Type:             enum.Cancel,
			MessageContextId: 3,
		},
	}

	var decoded GenericCallback
	roundTrip(original, &decoded, t)

	if decoded.Type != enum.Cancel {
		t.Errorf("Type = %v, want Cancel", decoded.Type)
	}
}

func TestTransactionTypeCallback_JSON(t *testing.T) {
	original := TransactionTypeCallback{
		Callback: Callback{
			Type:             enum.TransactionType,
			MessageContextId: 1,
		},
		TransactionTypeId: 5,
	}

	var decoded TransactionTypeCallback
	roundTrip(original, &decoded, t)

	if decoded.TransactionTypeId != 5 {
		t.Errorf("TransactionTypeId = %d, want 5", decoded.TransactionTypeId)
	}
}

func roundTrip(original interface{}, decoded interface{}, t *testing.T) {
	t.Helper()
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	err = json.Unmarshal(data, decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v (%s)", err, data)
	}
}
