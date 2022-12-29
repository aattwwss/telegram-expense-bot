package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
)

const PercentCategoryAmountMsg = "<code>%s%.1f%% %s %s%s\n</code>" // E.g. 82.8% Taxes    $1,234.00

type Transaction struct {
	Id          int
	Datetime    time.Time
	CategoryId  int
	Description string
	UserId      int64
	Amount      *money.Money
}

type Breakdown struct {
	CategoryName string
	Amount       *money.Money
	Percent      float64
}

type Breakdowns []Breakdown

func (bds Breakdowns) GetFormattedHTMLText() string {
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
