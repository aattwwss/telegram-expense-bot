package domain

import (
	"testing"
	"time"
)

func TestMonthlySummariesGetLongestLabelLength(t *testing.T) {
	s := MonthlySummaries{
		{TransactionTypeLabel: "Spent", Month: time.January, Year: 2023, Amount: 1000, Multiplier: -1},
		{TransactionTypeLabel: "Received", Month: time.January, Year: 2023, Amount: 5000, Multiplier: 1},
	}

	got := s.GetLongestLabelLength()
	want := len("Received") // 8
	if got != want {
		t.Errorf("GetLongestLabelLength() = %d, want %d", got, want)
	}
}

func TestMonthlySummariesGetLongestLabelLengthEmpty(t *testing.T) {
	s := MonthlySummaries{}
	got := s.GetLongestLabelLength()
	if got != 0 {
		t.Errorf("GetLongestLabelLength() on empty = %d, want 0", got)
	}
}

func TestMonthlySummariesGenerateReportText(t *testing.T) {
	s := MonthlySummaries{
		{
			Month:                time.January,
			Year:                 2023,
			Amount:               1000,
			TransactionTypeLabel: "Spent",
			Multiplier:           -1,
		},
		{
			Month:                time.January,
			Year:                 2023,
			Amount:               500,
			TransactionTypeLabel: "Received",
			Multiplier:           1,
		},
	}

	text := s.GenerateReportText("SGD")
	if len(text) == 0 {
		t.Error("expected non-empty report text")
	}
	if !contains(text, "Jan") || !contains(text, "2023") {
		t.Error("expected month/year header")
	}
	if !contains(text, "Spent") || !contains(text, "Received") {
		t.Error("expected transaction type labels")
	}
}

func TestMonthlySummaryGetPaddedSpacesForLabel(t *testing.T) {
	s := &MonthlySummary{TransactionTypeLabel: "Foo"}

	got := s.GetPaddedSpacesForLabel(6)
	want := "   " // 6 - 3 = 3 spaces
	if got != want {
		t.Errorf("GetPaddedSpacesForLabel(6) = %q (len=%d), want %q (len=%d)", got, len(got), want, len(want))
	}
}
