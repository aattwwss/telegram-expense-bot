package domain

import (
	"time"

	"github.com/Rhymond/go-money"
)

type Transaction struct {
	Id          int
	Datetime    time.Time
	CategoryId  int
	Description string
	UserId      int64
	Amount      *money.Money
}

type TransactionCategoryBreakdown struct {
	CategoryName string
	Amount       *money.Money
	Percent      float64
}
