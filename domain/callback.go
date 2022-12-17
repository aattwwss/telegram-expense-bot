package domain

type Callback struct {
	Type             string `json:"t,omitempty"`
	MessageContextId int    `json:"mc,omitempty"`
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
