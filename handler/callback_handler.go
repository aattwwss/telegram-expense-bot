package handler

import (
	"github.com/aattwwss/telegram-expense-bot/dao"
)

const ()

type CallbackHandler struct {
	userDao        dao.UserDAO
	transactionDao dao.TransactionDAO
}

func NewCallbackHandler(userDao dao.UserDAO, transactionDao dao.TransactionDAO) CallbackHandler {
	return CallbackHandler{
		userDao:        userDao,
		transactionDao: transactionDao,
	}
}
