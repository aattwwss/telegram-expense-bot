package entity

import "time"

type User struct {
	Id        int64
	IsBot     bool
	FirstName string
	LastName  *string
	Username  *string
}

type Transaction struct {
	Id          int
	Datetime    time.Time
	CategoryId  int
	Description string
	UserId      int64
	Amount      int64
	Currency    string
}
