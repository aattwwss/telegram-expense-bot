package message

const (
	HelpMsg = `
Message directly with the format "[amount] [description]" to start recording an expenses. 
Make sure the amount is purely a number or deciaml with no currency code or symbol.
The desciption can be made up of any characters.

Example:
✔️ "5.50 Chicken Rice" (without the quotes) to record an expense of $5.50 with the description "Chicken Rice".
✔️ "100 Ikea Table" (without the quotes) to record an expense of $100 with the description "Ikea Table".
❌ "Computer 2400" (without the quotes) will give an error".
❌ "$20.78 Pizza" (without the quotes) will give an error".
	
The recorded dollar ($) is the default currency symbol with support to up to 2 decimal places for the cents.

Type /stats [month] [year] to view the breakdown for the month.
Type /list [month] [year] to view the expenses for the month.
Type /export [month] [year] to export the expenses for the month.
Type /undo to revert the last recorded expenses.

List the expenses for current month and year
E.g. "/list".

List the expenses for the month of current year
E.g. "/list 2".
E.g. "/list Feb".
E.g. "/list February".

List the expenses for the month of February 2022
E.g. "/list 2 2022".
E.g. "/list Feb 2022".
E.g. "/list February 2022".

Stats and export follow the same rules as well!

If you have any questions or problems, email me at telegram.expense.bot@gmail.com
`

	TransactionTypeReplyMsg          = "Select a transaction type"
	TransactionStartReplyMsg         = "Select a category"
	TransactionEndReplyMsg           = "\n<i>%s</i>"
	TransactionLatestNotFound        = "You have no more transaction to delete."
	TransactionDeleteConfirmationMsg = "Do you want to delete your transaction of %s %s ?"
	TransactionDeletedReplyMsg       = "Your transaction of %s %s has been deleted."

	GenericErrReplyMsg = "Something went wrong :("
	WorkInProgressMsg  = "Sorry this function is still a work in progress."
)
