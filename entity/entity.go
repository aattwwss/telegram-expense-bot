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
	Id           int
	Datetime     time.Time
	CategoryId   int
	CategoryName string
	Description  string
	UserId       int64
	Amount       int64
	Currency     string
}

type Category struct {
	Id                int
	Name              string
	TransactionTypeId int
}

type MonthlySummary struct {
	Datetime             time.Time
	Amount               int64
	TransactionTypeLabel string
	Multiplier           int64
}

type TransactionType struct {
	Id         int
	Name       string
	Multiplier int
	ReplyText  string
}

type MessageContext struct {
	Id        int
	ChatId    int64
	MessageId int
	Message   string
	CreatedAt time.Time
}

type TransactionBreakdown struct {
	CategoryId   int
	CategoryName string
	Amount       int64
}
