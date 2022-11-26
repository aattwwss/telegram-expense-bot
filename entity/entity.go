package entity

import (
	"time"
)

type User struct {
	Id        int64
	IsBot     bool
	FirstName string
	LastName  *string
	Username  *string
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
	Date                 string // using string to represent the date aspect only
	Amount               int64
	TransactionTypeLabel string
	Multiplier           int
}
