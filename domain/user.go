package domain

import (
	"time"

	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/enum"
)

type User struct {
	Id             int64
	Locale         string
	Currency       *money.Currency
	Location       *time.Location
	CurrentContext enum.UserContext
}
