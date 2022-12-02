package domain

type Callback struct {
	TypeName string
}

type TransactionTypeCallback struct {
	Callback
	TransactionId int
}

type CategoryCallback struct {
	TransactionTypeCallback
	CategoryId int
}
