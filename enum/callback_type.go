package enum

import "strings"

type CallbackType string
type PaginateAction string
type UserContext string

const (
	TransactionType CallbackType = "TransactionType"
	Category        CallbackType = "Category"
	Pagination      CallbackType = "Pagination"
	Undo            CallbackType = "Undo"
	Cancel          CallbackType = "Cancel"

	Next     PaginateAction = "Next"
	Previous PaginateAction = "Prev"

	Transaction UserContext = "TRANSACTION"
	SetTimeZone UserContext = "SET_TIMEZONE"
	SetCurrency UserContext = "SET_CURRENCY"
)

var userContextMap = map[string]UserContext{
	"TRANSACTION":  Transaction,
	"SET_TIMEZONE": SetTimeZone,
	"SET_CURRENCY": SetCurrency,
}

func ParseUserContext(s string) (UserContext, bool) {
	uc, ok := userContextMap[strings.ToUpper(s)]
	return uc, ok
}
