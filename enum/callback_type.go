package enum

type CallbackType string
type PaginateAction string
type SortOrder string

const (
	TransactionType CallbackType = "TransactionType"
	Category        CallbackType = "Category"
	Pagination      CallbackType = "Pagination"
	Cancel          CallbackType = "Cancel"

	Next     PaginateAction = "Next"
	Previous PaginateAction = "Prev"

	ASC  SortOrder = "ASC"
	DESC SortOrder = "DESC"
)
