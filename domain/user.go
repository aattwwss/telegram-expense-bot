package domain

import (
	"time"

	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/entity"
)

type User struct {
	Id       int64
	Locale   string
	Currency *money.Currency
	Location *time.Location
}

func UserFromEntity(e entity.User) (*User, error) {
	loc, err := time.LoadLocation(e.Timezone)
	if err != nil {
		return nil, err
	}
	return &User{
		Id:       e.Id,
		Locale:   e.Locale,
		Currency: money.GetCurrency(e.Currency),
		Location: loc,
	}, nil
}
