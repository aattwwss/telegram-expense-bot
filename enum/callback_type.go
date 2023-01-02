package enum

type CallbackType string
type PaginateAction string

const (
	TransactionType CallbackType = "TransactionType"
	Category        CallbackType = "Category"
	Pagination      CallbackType = "Pagination"
	Cancel          CallbackType = "Cancel"

	Next     PaginateAction = "Next"
	Previous PaginateAction = "Prev"
)
