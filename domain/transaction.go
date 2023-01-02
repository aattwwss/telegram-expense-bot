package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
)

const PercentCategoryAmountMsg = "<code>%s%.1f%% %s %s%s\n</code>" // E.g. 82.8% Taxes    $1,234.00
const ListTransactionMsg = "<code>%s\n%s %s %s%s\n\n</code>"

type Transaction struct {
	Id           int
	Datetime     time.Time
	CategoryId   int
	CategoryName string
	Description  string
	UserId       int64
	Amount       *money.Money
}

type Transactions []Transaction

func (trxs Transactions) GetFormattedHTMLMsg() string {
	text := ""
	longest := 0
	for _, t := range trxs {
		length := len(t.CategoryName) + len(t.Description)
		if length > longest {
			longest = length
		}
	}
	// display the transactions in reverse order
	for i := len(trxs) - 1; i >= 0; i-- {
		t := trxs[i]
		dtString := t.Datetime.Format("02/01/06 15:04")
		spacesToPadAfterDesc := longest - len(t.CategoryName) - len(t.Description)
		text += fmt.Sprintf(ListTransactionMsg, dtString, t.CategoryName, t.Description, strings.Repeat(" ", spacesToPadAfterDesc), t.Amount.Display())
	}
	return text
}

type Breakdown struct {
	CategoryName string
	Amount       *money.Money
	Percent      float64
}

type Breakdowns []Breakdown

func (bds Breakdowns) GetFormattedHTMLMsg() string {
	text := ""

	longest := 0
	for _, b := range bds {
		length := len(b.CategoryName)
		if length > longest {
			longest = length
		}
	}

	for _, b := range bds {
		spacesToPadBeforePercent := ""
		if b.Percent < 10 {
			spacesToPadBeforePercent = " "
		}
		spacesToPadAfterCategory := longest - len(b.CategoryName)
		text += fmt.Sprintf(PercentCategoryAmountMsg, spacesToPadBeforePercent, b.Percent, b.CategoryName, strings.Repeat(" ", spacesToPadAfterCategory), b.Amount.Display())
	}
	return text
}
