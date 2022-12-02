package handler

import (
	"context"
	"fmt"
	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/message"
	"github.com/aattwwss/telegram-expense-bot/repo"
	"github.com/aattwwss/telegram-expense-bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

const (
	categoriesInlineColSize = 3
)

type CallbackHandler struct {
	userRepo            repo.UserRepo
	transactionRepo     repo.TransactionRepo
	transactionTypeRepo repo.TransactionTypeRepo
	categoryRepo        repo.CategoryRepo
}

func NewCallbackHandler(userRepo repo.UserRepo, transactionRepo repo.TransactionRepo, transactionTypeRepo repo.TransactionTypeRepo, categoryRepo repo.CategoryRepo) CallbackHandler {
	return CallbackHandler{
		userRepo:            userRepo,
		transactionRepo:     transactionRepo,
		transactionTypeRepo: transactionTypeRepo,
		categoryRepo:        categoryRepo,
	}
}

func (handler CallbackHandler) FromTransactionType(ctx context.Context, msg *tgbotapi.MessageConfig, callbackQuery *tgbotapi.CallbackQuery, data string) {
	text := callbackQuery.Message.Text

	transactionTypeId, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		log.Error().Msgf("FromTransactionType error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	categories, err := handler.categoryRepo.FindByTransactionTypeId(ctx, transactionTypeId)
	if err != nil {
		log.Error().Msgf("%V", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	moneyAmount, err := parseMoneyFromCallback(text, message.TransactionTypeReplyMsg, money.SGD)
	if err != nil {
		log.Error().Msgf("%V", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	msg.Text = fmt.Sprintf(message.TransactionReplyMsg+"%v", moneyAmount.AsMajorUnits())
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: newCategoriesKeyboard(categories, categoriesInlineColSize)}

}

func (handler CallbackHandler) FromCategory(ctx context.Context, msg *tgbotapi.MessageConfig, callbackQuery *tgbotapi.CallbackQuery, data string) {
	text := callbackQuery.Message.Text

	categoryId, err := strconv.Atoi(data)
	if err != nil {
		log.Error().Msgf("FromCategory error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}
	category, err := handler.categoryRepo.GetById(ctx, categoryId)
	if err != nil {
		msg.Text = message.GenericErrReplyMsg
		return
	}
	moneyTransacted, err := parseMoneyFromCallback(text, message.TransactionReplyMsg, money.SGD)
	if err != nil {
		log.Error().Msgf("FromCategory error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	transaction := domain.Transaction{
		Datetime:    time.Now(),
		CategoryId:  categoryId,
		Description: "",
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

func parseMoneyFromCallback(s string, msg string, currencyCode string) (*money.Money, error) {
	floatString := strings.ReplaceAll(s, msg, "")
	floatAmount, err := strconv.ParseFloat(floatString, 10)
	if err != nil {
		return nil, err
	}
	return money.NewFromFloat(floatAmount, currencyCode), nil
}

func newCategoriesKeyboard(categories []domain.Category, colSize int) [][]tgbotapi.InlineKeyboardButton {
	var configs []util.InlineKeyboardConfig
	for _, category := range categories {
		config := util.NewInlineKeyboardConfig(category.Name, util.CallbackDataSerialize(category, category.Id))
		configs = append(configs, config)
	}

	return util.NewInlineKeyboard(configs, colSize)

}
