package domain

import "github.com/aattwwss/telegram-expense-bot/enum"

type Callback struct {
	Type             enum.CallbackType `json:"t,omitempty"`
	MessageContextId int               `json:"mc,omitempty"`
}

type GenericCallback struct {
	Callback `json:"c"`
}

type TransactionTypeCallback struct {
	Callback          `json:"c"`
	TransactionTypeId int `json:"id,omitempty"`
}

type CategoryCallback struct {
	Callback   `json:"c"`
	CategoryId int `json:"id,omitempty"`
}

type PaginationCallback struct {
	Callback `json:"c"`
	Action   enum.PaginateAction `json:"a,omitempty"`
	Offset   int                 `json:"o"`
	Limit    int                 `json:"l"`
}

type UndoCallback struct {
	Callback      `json:"c"`
	TransactionId int `json:"t"`
}
