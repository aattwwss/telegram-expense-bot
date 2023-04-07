package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
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

func (handler CallbackHandler) FromCategory(ctx context.Context, bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	user, err := handler.userRepo.FindUserById(ctx, callbackQuery.From.ID)
	if err != nil {
		log.Error().Msgf("Error finding user for category: %v", err)
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	var categoryCallback domain.CategoryCallback
	err = json.Unmarshal([]byte(callbackQuery.Data), &categoryCallback)
	if err != nil {
		log.Error().Msgf("FromCategory unmarshall error: %v", err)
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	defer handler.deleteMessageContext(ctx, categoryCallback.MessageContextId)

	category, err := handler.categoryRepo.GetById(ctx, categoryCallback.CategoryId)
	if err != nil {
		log.Error().Msgf("Get category by id error: %v", err)
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	messageContext, err := handler.messageContextRepo.GetMessageById(ctx, categoryCallback.Callback.MessageContextId)
	if err != nil {
		log.Error().Msgf("Get message context by id error: %v", err)
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	amountString, err := parseFloatStringFromString(messageContext)
	if err != nil {
		log.Error().Msgf("Parsing float string from meesage context error: %v", err)
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	amountFloat, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	amountInt, err := decimalise(amountFloat, *user.Currency)
	if err != nil {
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}
	stringAfter := util.After(messageContext, amountString)
	description := strings.TrimSpace(stringAfter)

	moneyTransacted := money.New(amountInt, user.Currency.Code)

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
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	transactionType, err := handler.transactionTypeRepo.GetById(ctx, category.TransactionTypeId)
	if err != nil {
		log.Error().Msgf("FromCategory error: %v", err)
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	text := fmt.Sprintf(transactionType.ReplyText, moneyTransacted.Display(), category.Name)
	text += fmt.Sprintf(message.TransactionEndReplyMsg, description)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	util.BotSendWrapper(bot, msg)
}

func (handler CallbackHandler) FromPagination(ctx context.Context, bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	// TODO Find a way to handle the persisting context when paginating
	userId := callbackQuery.From.ID
	user, err := handler.userRepo.FindUserById(ctx, userId)
	if err != nil {
		log.Error().Msgf("Error finding user for stats: %w", err)
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	var paginationCallback domain.PaginationCallback
	err = json.Unmarshal([]byte(callbackQuery.Data), &paginationCallback)
	if err != nil {
		log.Error().Msgf("FromPagination unmarshall error: %w", err)
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	messageContext, err := handler.messageContextRepo.GetMessageById(ctx, paginationCallback.Callback.MessageContextId)
	if err != nil {
		log.Error().Msgf("Get message context by id error: %w", err)
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	month, year := util.ParseMonthYearFromMessage(messageContext)

	offset, limit := paginationCallback.Offset, paginationCallback.Limit
	transactions, totalCount, err := handler.transactionRepo.ListByMonthAndYear(ctx, month, year, offset, limit, false, *user)

	inlineKeyboard, err := util.NewPaginationKeyboard(totalCount, offset, limit, paginationCallback.MessageContextId, 2)
	if err != nil {
		log.Error().Msgf("Error generating keyboard for transaction pagination: %w", err)
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	text := transactions.GetFormattedHTMLMsg(month, year, user.Location, totalCount, offset, limit)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, text)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}
	msg.ParseMode = tgbotapi.ModeHTML
	util.BotSendWrapper(bot, msg)
}

func (handler CallbackHandler) FromUndo(ctx context.Context, bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	userId := callbackQuery.From.ID
	var undoCallback domain.UndoCallback

	err := json.Unmarshal([]byte(callbackQuery.Data), &undoCallback)
	if err != nil {
		log.Error().Msgf("FromUndo unmarshall error: %w", err)
		return
	}
	log.Info().Msgf("transaction: %v", undoCallback.TransactionId)

	transaction, err := handler.transactionRepo.GetById(ctx, undoCallback.TransactionId, userId)
	if err != nil {
		log.Error().Msgf("FromUndo cannot find transaction error: %w", err)
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	err = handler.transactionRepo.DeleteById(ctx, undoCallback.TransactionId, userId)
	if err != nil {
		log.Error().Msgf("Error deleting latest transaction: %w", err)
		util.BotSendMessage(bot, callbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	text := fmt.Sprintf(message.TransactionDeletedReplyMsg, transaction.Amount.Display(), transaction.Description)
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, text)
	util.BotSendWrapper(bot, msg)
}

func (handler CallbackHandler) FromCancel(ctx context.Context, bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
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

var floatParser = regexp.MustCompile(`-?\d[\d,]*[.]?[\d{2}]*`)

func parseFloatStringFromString(s string) (string, error) {
	matches := floatParser.FindAllString(s, -1)
	if len(matches) == 0 {
		return "", errors.New("no float found in string: " + s)
	}
	return matches[0], nil
}
func (handler CallbackHandler) deleteMessageContext(ctx context.Context, id int) {
	err := handler.messageContextRepo.DeleteById(ctx, id)
	if err != nil {
		log.Error().Msgf("deleteMessageContext error: %w", err)
	}
}

// decimalise decimalise the value of a currency to its lowest denomination
func decimalise(value float64, currency money.Currency) (int64, error) {
	formatString := fmt.Sprintf("%%.%df", currency.Fraction)
	formatted := fmt.Sprintf(formatString, value)
	intString := removeNonNumeric(formatted)
	res, err := strconv.ParseInt(intString, 10, 64)
	if err != nil {
		return 0, err
	}
	return res, err
}

// removeNonNumeric removes all non-numeric characters from a string
func removeNonNumeric(s string) string {
	var sb strings.Builder
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			sb.WriteRune(ch)
		}
	}
	return sb.String()
}
