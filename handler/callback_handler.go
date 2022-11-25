package handler

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/dao"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type CallbackHandler struct {
	userDao        dao.UserDAO
	transactionDao dao.TransactionDAO
	categoryDao    dao.CategoryDAO
}

func NewCallbackHandler(userDao dao.UserDAO, transactionDao dao.TransactionDAO, categoryDao dao.CategoryDAO) CallbackHandler {
	return CallbackHandler{
		userDao:        userDao,
		transactionDao: transactionDao,
		categoryDao:    categoryDao,
	}
}

func (handler CallbackHandler) FromCategory(ctx context.Context, msg *tgbotapi.MessageConfig, callbackQuery *tgbotapi.CallbackQuery, idString string) {

	id, err := strconv.Atoi(idString)
	if err != nil {
		msg.Text = "Something went wrong :("
		return
	}
	category, err := handler.categoryDao.GetById(ctx, id)
	if err != nil {
		msg.Text = "Something went wrong :("
		return
	}
	msg.Text = category.Name
}
