package handler

import (
	"context"
	"fmt"
	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/entity"
	"github.com/aattwwss/telegram-expense-bot/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
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

//func (handler CallbackHandler) FromTransactionType(ctx context.Context, msg *tgbotapi.MessageConfig, callbackQuery *tgbotapi.CallbackQuery, data string) {
//	text := callbackQuery.Message.Text
//
//	transactionTypeId, err := strconv.Atoi(data)
//	if err != nil {
//		log.Error().Msgf("FromTransactionType error: %v", err)
//		msg.Text = message.GenericErrReplyMsg
//		return
//	}
//
//}

func (handler CallbackHandler) FromCategory(ctx context.Context, msg *tgbotapi.MessageConfig, callbackQuery *tgbotapi.CallbackQuery, data string) {
	text := callbackQuery.Message.Text

	categoryId, err := strconv.Atoi(data)
	if err != nil {
		log.Error().Msgf("FromCategory error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}
	category, err := handler.categoryDao.GetById(ctx, categoryId)
	if err != nil {
		msg.Text = message.GenericErrReplyMsg
		return
	}
	moneyTransacted, err := parseMoneyFromTransactionCallback(text, money.SGD)
	if err != nil {
		log.Error().Msgf("FromCategory error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	transaction := entity.Transaction{
		Datetime:    time.Now(),
		CategoryId:  categoryId,
		Description: "",
		UserId:      callbackQuery.From.ID,
		Amount:      moneyTransacted.Amount(),
		Currency:    money.SGD,
	}

	err = handler.transactionDao.Insert(ctx, transaction)
	if err != nil {
		log.Error().Msgf("FromCategory error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}
	msg.Text = fmt.Sprintf("You spent %s on %s", moneyTransacted.Display(), category.Name)
}

func parseMoneyFromTransactionCallback(s string, currencyCode string) (*money.Money, error) {
	floatString := strings.ReplaceAll(s, message.TransactionReplyMsg, "")
	floatAmount, err := strconv.ParseFloat(floatString, 10)
	if err != nil {
		return nil, err
	}
	return money.NewFromFloat(floatAmount, currencyCode), nil
}
