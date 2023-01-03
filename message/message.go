package message

const (
	HelpMsg = `Type /stats [month] [year] to view the breakdown for the month.
Type /list [month] [year] to view the expenses for the month.
Type /undo to revert the last recorded expenses.

Message directly to start recording an expenses.
E.g. "5.50 Chicken Rice" to record an expense of $5.50 with the description "Chicken Rice".`

	TransactionTypeReplyMsg    = "Select a transaction type"
	TransactionStartReplyMsg   = "Select a category"
	TransactionEndReplyMsg     = "\n<i>%s</i>"
	TransactionLatestNotFound  = "You have no more transaction to delete."
	TransactionDeletedReplyMsg = "Your transaction of %s %s has been deleted."

	GenericErrReplyMsg = "Something went wrong :("
	WorkInProgressMsg  = "Sorry this function is still a work in progress."
)
