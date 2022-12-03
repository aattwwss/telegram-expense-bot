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
	Datetime             time.Time
	Amount               int64
	TransactionTypeLabel string
	Multiplier           int64
}

type TransactionType struct {
	Id         int64
	Name       string
	Multiplier int
	ReplyText  string
}
