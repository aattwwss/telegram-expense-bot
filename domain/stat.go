package domain

import (
	"fmt"
	"github.com/Rhymond/go-money"
	"strings"
	"time"
)

const (
	monthYearHeaderHTMLMsg    = "<code>\n%v %v\n</code>"
	transactionSummaryHTMLMsg = "<code>%v:%s %v\n</code>"
	transactionTotalHTMLMsg   = "<code>🟡 Total: %v\n</code>"
)

type MonthlySummaries []MonthlySummary

func (s MonthlySummaries) GetLongestLabelLength() int {
	longestLabel := 0
	for _, summary := range s {
		lengthOfLabel := len(summary.TransactionTypeLabel)
		if lengthOfLabel > longestLabel {
			longestLabel = lengthOfLabel
		}
	}
	return longestLabel
}

func (s MonthlySummaries) GenerateReportText() string {
	currMonth := ""
	var totalAmountForTheMonth int64
	var msg string

	longestLabel := s.GetLongestLabelLength()

	for i, summary := range s {
		month := summary.Month.String()[:3]
		if currMonth != month {
			msg += fmt.Sprintf(monthYearHeaderHTMLMsg, month, summary.Year)
			currMonth = month
			totalAmountForTheMonth = 0
		}

		totalAmountForTheMonth += summary.Amount * summary.Multiplier
		moneyAmount := money.New(summary.Amount, money.SGD)
		msg += fmt.Sprintf(transactionSummaryHTMLMsg, summary.TransactionTypeLabel, summary.GetPaddedSpacesForLabel(longestLabel), moneyAmount.Display())

		if i == len(s)-1 || s[i+1].Month.String()[:3] != currMonth {
			msg += fmt.Sprintf(transactionTotalHTMLMsg, money.New(totalAmountForTheMonth, money.SGD).Display())
		}
	}
	return msg
}

type MonthlySummary struct {
	Month                time.Month
	Year                 int
	Amount               int64
	TransactionTypeLabel string
	Multiplier           int64
}

func (s *MonthlySummary) GetPaddedSpacesForLabel(lengthToPadTo int) string {
	return strings.Repeat(" ", lengthToPadTo-len(s.TransactionTypeLabel))
}
