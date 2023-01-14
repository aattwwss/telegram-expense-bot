package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/enum"
	"github.com/aattwwss/telegram-expense-bot/message"
	"github.com/aattwwss/telegram-expense-bot/repo"
	"github.com/aattwwss/telegram-expense-bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
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

func (handler CallbackHandler) FromCategory(ctx context.Context, msg *tgbotapi.MessageConfig, callbackQuery *tgbotapi.CallbackQuery) {
	var categoryCallback domain.CategoryCallback
	err := json.Unmarshal([]byte(callbackQuery.Data), &categoryCallback)
	if err != nil {
		log.Error().Msgf("FromCategory unmarshall error: %w", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	defer handler.deleteMessageContext(ctx, categoryCallback.MessageContextId)

	category, err := handler.categoryRepo.GetById(ctx, categoryCallback.CategoryId)
	if err != nil {
		log.Error().Msgf("Get category by id error: %w", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	messageContext, err := handler.messageContextRepo.GetMessageById(ctx, categoryCallback.Callback.MessageContextId)
	if err != nil {
		log.Error().Msgf("Get message context by id error: %w", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	amountString, err := util.ParseFloatStringFromString(messageContext)
	if err != nil {
		log.Error().Msgf("%w", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	amountFloat, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		msg.Text = message.GenericErrReplyMsg
		return
	}

	stringAfter := util.After(messageContext, amountString)
	description := strings.TrimSpace(stringAfter)

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
		log.Error().Msgf("FromCategory error: %w", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	transactionType, err := handler.transactionTypeRepo.GetById(ctx, category.TransactionTypeId)
	if err != nil {
		log.Error().Msgf("FromCategory error: %w", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	text := fmt.Sprintf(transactionType.ReplyText, moneyTransacted.Display(), category.Name) + "\n"
	text += fmt.Sprintf(message.TransactionEndReplyMsg, description)
	msg.Text = text
	msg.ParseMode = tgbotapi.ModeHTML
}

func (handler CallbackHandler) FromPagination(ctx context.Context, msg *tgbotapi.MessageConfig, callbackQuery *tgbotapi.CallbackQuery) {
	// TODO Find a way to handle the peristing context when paginating
	userId := callbackQuery.From.ID
	user, err := handler.userRepo.FindUserById(ctx, userId)
	if err != nil {
		log.Error().Msgf("Error finding user for stats: %w", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	var paginationCallback domain.PaginationCallback
	err = json.Unmarshal([]byte(callbackQuery.Data), &paginationCallback)
	if err != nil {
		log.Error().Msgf("FromPagination unmarshall error: %w", err)
		return
	}

	messageContext, err := handler.messageContextRepo.GetMessageById(ctx, paginationCallback.Callback.MessageContextId)
	if err != nil {
		log.Error().Msgf("Get message context by id error: %w", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	month, year := util.ParseMonthYearFromMessage(messageContext)

	offset, limit := paginationCallback.Offset, paginationCallback.Limit
	transactions, totalCount, err := handler.transactionRepo.ListByMonthAndYear(ctx, month, year, offset, limit, *user)

	inlineKeyboard, err := util.NewPaginationKeyboard(totalCount, offset, limit, paginationCallback.MessageContextId, 2)
	if err != nil {
		log.Error().Msgf("Error generating keyboard for transaction pagination: %w", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}

	msg.Text = transactions.GetFormattedHTMLMsg(month, year, user.Location, totalCount, offset, limit)
	msg.ParseMode = tgbotapi.ModeHTML
}

func (handler CallbackHandler) FromCancel(ctx context.Context, msg *tgbotapi.MessageConfig, callbackQuery *tgbotapi.CallbackQuery) {

	var genericCallback domain.GenericCallback
	err := json.Unmarshal([]byte(callbackQuery.Data), &genericCallback)
	if err != nil {
		log.Error().Msgf("FromCancel unmarshall error: %w", err)
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
		log.Error().Msgf("deleteMessageContext error: %w", err)
	}
}
