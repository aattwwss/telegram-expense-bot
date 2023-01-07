package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
)

const PercentCategoryAmountMsg = "<code>%s%.1f%% %s %s%s\n</code>" // E.g. 82.8% Taxes    $1,234.00
const ListTransactionHeader = "<b>%s %v</b>\n\n"                   // E.g. January 2023
const ListTransactionBody = "<code>%s\n%s %s %s%s\n\n</code>"
const ListTransactionFooter = "<code>[%v/%v]</code>" //E.g. [1/3]

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

func (trxs Transactions) GetFormattedHTMLMsg(searchedMonth time.Month, searchedYear int, loc *time.Location, totalCount int, currentOffset int, pageSize int) string {
	text := fmt.Sprintf(ListTransactionHeader, searchedMonth.String(), searchedYear)
	longest := 0

	for _, t := range trxs {
		length := len(t.CategoryName) + len(t.Description)
		if length > longest {
			longest = length
		}
	}

	for _, t := range trxs {
		dtString := t.Datetime.In(loc).Format("02/01/06 15:04")
		spacesToPadAfterDesc := longest - len(t.CategoryName) - len(t.Description)
		text += fmt.Sprintf(ListTransactionBody, dtString, t.CategoryName, t.Description, strings.Repeat(" ", spacesToPadAfterDesc), t.Amount.Display())
	}

	numOfPages := (totalCount-1)/pageSize + 1
	currentPage := (currentOffset)/pageSize + 1
	text += fmt.Sprintf(ListTransactionFooter, currentPage, numOfPages)
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
