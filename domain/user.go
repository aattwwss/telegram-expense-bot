package domain

import (
	"github.com/Rhymond/go-money"
	"time"
)

type User struct {
	Id       int64
	Locale   string
	Currency *money.Currency
	Location *time.Location
}
