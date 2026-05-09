package entity

import "time"

type TransactionListQuery struct {
	Month    time.Month
	Year     int
	Offset   int
	Limit    int
	Asc      bool
	UserId   int64
	Location *time.Location
}
