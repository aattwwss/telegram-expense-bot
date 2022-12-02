package entity

import (
	"time"
)

type User struct {
	Id       int64
	Locale   string
	Currency string
	Timezone string
}

type Transaction struct {
	Id          int64
	Datetime    time.Time
	CategoryId  int
	Description string
	UserId      int64
	Amount      int64
	Currency    string
}

type Category struct {
	Id                int64
	Name              string
	TransactionTypeId int64
}

type MonthlySummary struct {
	Month                time.Month
	Year                 int
	Amount               int64
	TransactionTypeLabel string
	Multiplier           int64

	spacesToPad int
}

func (s *MonthlySummary) SetSpacesToPad(n int) {
	s.spacesToPad = n
}

func (s *MonthlySummary) GetSpacesToPad() int {
	return s.spacesToPad
}
