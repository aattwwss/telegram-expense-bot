package domain

import (
	"github.com/Rhymond/go-money"
	"time"
)

type Transaction struct {
	Id          int
	Datetime    time.Time
	CategoryId  int
	Description string
	UserId      int64
	Amount      *money.Money
}
