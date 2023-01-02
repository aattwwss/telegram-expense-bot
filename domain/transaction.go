package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
)

const PercentCategoryAmountMsg = "<code>%s%.1f%% %s %s%s\n</code>" // E.g. 82.8% Taxes    $1,234.00
// const ListTransactionMsg = "<code>%s\n%s %s\n%s\n\n</code>"        // E.g. 12/01/2022 Food hotdog $123.45
const ListTransactionMsg = "%s\n<b>%s</b> %s - %s\n\n"

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
	// display the transactions in reverse order
	for i := len(trxs) - 1; i >= 0; i-- {
		t := trxs[i]
		text += fmt.Sprintf(ListTransactionMsg, t.Datetime.Format("<b>02/01/06</b> 15:04"), t.CategoryName, t.Description, t.Amount.Display())
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
