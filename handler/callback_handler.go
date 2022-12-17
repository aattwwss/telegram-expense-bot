package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/enum"
	"github.com/aattwwss/telegram-expense-bot/message"
	"github.com/aattwwss/telegram-expense-bot/repo"
	"github.com/aattwwss/telegram-expense-bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

const (
	categoriesInlineColSize = 3
)

type CallbackHandler struct {
	userRepo            repo.UserRepo
	transactionRepo     repo.TransactionRepo
	messageContextRepo  repo.MessageContextRepo
	transactionTypeRepo repo.TransactionTypeRepo
	categoryRepo        repo.CategoryRepo
}

func NewCallbackHandler(userRepo repo.UserRepo, transactionRepo repo.TransactionRepo, messageContextRepo repo.MessageContextRepo, transactionTypeRepo repo.TransactionTypeRepo, categoryRepo repo.CategoryRepo) CallbackHandler {
	return CallbackHandler{
		userRepo:            userRepo,
		transactionRepo:     transactionRepo,
		messageContextRepo:  messageContextRepo,
		transactionTypeRepo: transactionTypeRepo,
		categoryRepo:        categoryRepo,
	}
}

func (handler CallbackHandler) FromTransactionType(ctx context.Context, msg *tgbotapi.MessageConfig, callbackQuery *tgbotapi.CallbackQuery) {

	var transactionTypeCallback domain.TransactionTypeCallback
	err := json.Unmarshal([]byte(callbackQuery.Data), &transactionTypeCallback)
	if err != nil {
		log.Error().Msgf("FromTransactionType unmarshall error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	categories, err := handler.categoryRepo.FindByTransactionTypeId(ctx, transactionTypeCallback.TransactionTypeId)
	if err != nil {
		log.Error().Msgf("FindByTransactionTypeId error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	inlineKeyboard, err := newCategoriesKeyboard(categories, transactionTypeCallback.Callback.MessageContextId, categoriesInlineColSize)
	if err != nil {
		log.Error().Msgf("newCategoriesKeyboard error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	msg.Text = message.TransactionReplyMsg
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}

}

func (handler CallbackHandler) FromCategory(ctx context.Context, msg *tgbotapi.MessageConfig, callbackQuery *tgbotapi.CallbackQuery) {
	var categoryCallback domain.CategoryCallback
	err := json.Unmarshal([]byte(callbackQuery.Data), &categoryCallback)
	if err != nil {
		log.Error().Msgf("FromCategory unmarshall error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	defer handler.deleteMessageContext(ctx, categoryCallback.MessageContextId)

	category, err := handler.categoryRepo.GetById(ctx, categoryCallback.CategoryId)
	if err != nil {
		log.Error().Msgf("Get category by id error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	messageContext, err := handler.messageContextRepo.GetMessageById(ctx, categoryCallback.Callback.MessageContextId)
	if err != nil {
		log.Error().Msgf("Get message context by id error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	amountString, err := util.ParseFloatStringFromString(messageContext)
	if err != nil {
		log.Error().Msgf("%v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	amountFloat, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		msg.Text = message.GenericErrReplyMsg
		return
	}

	description := util.After(messageContext, amountString)

	moneyTransacted := money.NewFromFloat(amountFloat, money.SGD)

	transaction := domain.Transaction{
		Datetime:    time.Now(),
		CategoryId:  category.Id,
		Description: description,
		UserId:      callbackQuery.From.ID,
		Amount:      moneyTransacted,
	}

	err = handler.transactionRepo.Add(ctx, transaction)
	if err != nil {
		log.Error().Msgf("FromCategory error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	transactionType, err := handler.transactionTypeRepo.GetById(ctx, category.TransactionTypeId)
	if err != nil {
		log.Error().Msgf("FromCategory error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	msg.Text = fmt.Sprintf(transactionType.ReplyText, moneyTransacted.Display(), category.Name)
}

func (handler CallbackHandler) FromCancel(ctx context.Context, msg *tgbotapi.MessageConfig, callbackQuery *tgbotapi.CallbackQuery) {

	var genericCallback domain.GenericCallback
	err := json.Unmarshal([]byte(callbackQuery.Data), &genericCallback)
	if err != nil {
		log.Error().Msgf("FromCancel unmarshall error: %v", err)
		return
	}

	handler.deleteMessageContext(ctx, genericCallback.MessageContextId)
}

func newCategoriesKeyboard(categories []domain.Category, messageContextId int, colSize int) ([][]tgbotapi.InlineKeyboardButton, error) {
	var configs []util.InlineKeyboardConfig
	for _, category := range categories {
		data := domain.CategoryCallback{
			Callback: domain.Callback{
				Type:             enum.Category,
				MessageContextId: messageContextId,
			},
			CategoryId: category.Id,
		}

		dataJson, err := util.ToJson(data)
		if err != nil {
			return nil, err
		}

		config := util.NewInlineKeyboardConfig(category.Name, dataJson)
		configs = append(configs, config)
	}

	return util.NewInlineKeyboard(configs, messageContextId, colSize, true), nil

}

func (handler CallbackHandler) deleteMessageContext(ctx context.Context, id int) {
	err := handler.messageContextRepo.DeleteById(ctx, id)
	if err != nil {
		log.Error().Msgf("deleteMessageContext error: %v", err)
	}
}
