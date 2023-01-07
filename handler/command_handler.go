package handler

import (
	"context"
	"fmt"
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
	userExistsMsg            = "Welcome back! These are the summary of your transactions: \n"
	errorFindingUserMsg      = "Sorry there is a problem fetching your information.\n"
	errorCreatingUserMsg     = "Sorry there is a problem signing you up.\n"
	signUpSuccessMsg         = "Congratulations! We can get you started right away!\n"
	cannotRecogniseAmountMsg = "I don't recognise that amount of money :(\n"
	descriptionTooLong       = "Sorry, your description (max 20 characters) is too long :( \n"
	transactionListEmptyMsg  = "You have no transactions this month."

	statsHeaderHTMLMsg = "<b>%s %v\n</b>%s\n\n" // E.g. November 2022

	transactionTypeInlineColSize = 2

	transactionListDefaultPageSize = 10
)

type CommandHandler struct {
	transactionRepo     repo.TransactionRepo
	messageContextRepo  repo.MessageContextRepo
	transactionTypeRepo repo.TransactionTypeRepo
	categoryRepo        repo.CategoryRepo
	statRepo            repo.StatRepo
	userRepo            repo.UserRepo
}

func NewCommandHandler(userRepo repo.UserRepo, transactionRepo repo.TransactionRepo, messageContextRepo repo.MessageContextRepo, transactionTypeRepo repo.TransactionTypeRepo, categoryRepo repo.CategoryRepo, statRepo repo.StatRepo) CommandHandler {
	return CommandHandler{
		userRepo:            userRepo,
		transactionRepo:     transactionRepo,
		messageContextRepo:  messageContextRepo,
		transactionTypeRepo: transactionTypeRepo,
		categoryRepo:        categoryRepo,
		statRepo:            statRepo,
	}
}

func (handler CommandHandler) Start(ctx context.Context, msg *tgbotapi.MessageConfig, update tgbotapi.Update) {
	msg.Text = "Welcome to your expense tracker!\n"

	teleUser := update.SentFrom()

	dbUser, err := handler.userRepo.FindUserById(ctx, teleUser.ID)
	if err != nil {
		log.Error().Msgf("error finding user: %v", err)
		msg.Text += errorFindingUserMsg
		return
	}

	if dbUser != nil {
		log.Info().Msgf("User already exists. id: %v", dbUser.Id)
		msg.Text += userExistsMsg
		return
	}

	defaultLocation, _ := time.LoadLocation("Asia/Singapore")
	defaultCurrency := money.GetCurrency(money.SGD)

	user := domain.User{
		Id:       teleUser.ID,
		Locale:   "en",
		Currency: defaultCurrency,
		Location: defaultLocation,
	}

	err = handler.userRepo.Add(ctx, user)
	if err != nil {
		log.Error().Msgf("error adding user: %v", err)
		msg.Text += errorCreatingUserMsg
		return
	}
	msg.Text += signUpSuccessMsg
	return
}

func (handler CommandHandler) Help(ctx context.Context, msg *tgbotapi.MessageConfig, update tgbotapi.Update) {
	msg.Text = message.HelpMsg
	return
}

func (handler CommandHandler) Undo(ctx context.Context, msg *tgbotapi.MessageConfig, update tgbotapi.Update) {
	userId := update.Message.From.ID
	latestTransaction, err := handler.transactionRepo.FindLastestByUserId(ctx, userId)
	if err != nil {
		log.Error().Msgf("Error finding latest transaction: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}
	if latestTransaction == nil {
		msg.Text = message.TransactionLatestNotFound
		return
	}
	err = handler.transactionRepo.DeleteById(ctx, latestTransaction.Id, userId)
	if err != nil {
		log.Error().Msgf("Error deleting latest transaction: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}
	msg.Text = fmt.Sprintf(message.TransactionDeletedReplyMsg, latestTransaction.Amount.Display(), latestTransaction.Description)
	return
}

func (handler CommandHandler) StartTransaction(ctx context.Context, msg *tgbotapi.MessageConfig, update tgbotapi.Update) {
	userId := update.SentFrom().ID
	user, err := handler.userRepo.FindUserById(ctx, userId)
	if err != nil {
		log.Error().Msgf("Error finding user for transact: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}
	if user == nil {
		log.Error().Msgf("User not found for transact: %v", userId)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	floatString, err := util.ParseFloatStringFromString(update.Message.Text)
	if err != nil {
		log.Error().Msgf("%v", err)
		msg.Text = cannotRecogniseAmountMsg
		return
	}

	stringAfter := util.After(update.Message.Text, floatString)
	if len(strings.TrimSpace(stringAfter)) > 20 {
		msg.Text = descriptionTooLong
		return
	}

	contextId, err := handler.messageContextRepo.Add(ctx, update.Message.Chat.ID, update.Message.MessageID, update.Message.Text)
	if err != nil {
		log.Error().Msgf("Add message context error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	categories, err := handler.categoryRepo.FindAll(ctx)
	if err != nil {
		log.Error().Msgf("FindAll categories error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	inlineKeyboard, err := newCategoriesKeyboard(categories, contextId, categoriesInlineColSize)
	if err != nil {
		log.Error().Msgf("newCategoriesKeyboard error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	msg.Text = message.TransactionTypeReplyMsg
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}

}

func (handler CommandHandler) Stats(ctx context.Context, msg *tgbotapi.MessageConfig, update tgbotapi.Update) {
	userId := update.SentFrom().ID
	user, err := handler.userRepo.FindUserById(ctx, userId)
	if err != nil {
		log.Error().Msgf("Error finding user for stats: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	month, year := util.ParseMonthYearFromMessage(update.Message.Text)

	breakdowns, total, err := handler.transactionRepo.GetTransactionBreakdownByCategory(ctx, month, year, *user)

	if err != nil {
		log.Error().Msgf("Error getting breakdowns: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	header := fmt.Sprintf(statsHeaderHTMLMsg, month.String(), year, total.Display())
	msg.Text = header + breakdowns.GetFormattedHTMLMsg()
	msg.ParseMode = tgbotapi.ModeHTML
	return
}

func (handler CommandHandler) List(ctx context.Context, msg *tgbotapi.MessageConfig, update tgbotapi.Update) {
	pageSize := transactionListDefaultPageSize

	userId := update.SentFrom().ID
	user, err := handler.userRepo.FindUserById(ctx, userId)
	if err != nil {
		log.Error().Msgf("Error finding user for stats: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	contextId, err := handler.messageContextRepo.Add(ctx, update.Message.Chat.ID, update.Message.MessageID, update.Message.Text)
	if err != nil {
		log.Error().Msgf("Add message context error: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	month, year := util.ParseMonthYearFromMessage(update.Message.Text)

	transactions, totalCount, err := handler.transactionRepo.ListByMonthAndYear(ctx, month, year, 0, pageSize, *user)
	if err != nil {
		log.Error().Msgf("Error getting list of transactions: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}
	if totalCount == 0 {
		msg.Text = transactionListEmptyMsg
		return
	}

	inlineKeyboard, err := util.NewPaginationKeyboard(totalCount, 0, pageSize, contextId, 2)
	if err != nil {
		log.Error().Msgf("Error generating keyboard for transaction pagination: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}

	msg.Text = transactions.GetFormattedHTMLMsg(month, year, user.Location, totalCount, 0, pageSize)
	msg.ParseMode = tgbotapi.ModeHTML
	return
}

func newTransactionTypesKeyboard(transactionTypes []domain.TransactionType, messageContextId int, colSize int) ([][]tgbotapi.InlineKeyboardButton, error) {
	var configs []util.InlineKeyboardConfig
	for _, transactionType := range transactionTypes {
		data := domain.TransactionTypeCallback{
			Callback: domain.Callback{
				Type:             enum.TransactionType,
				MessageContextId: messageContextId,
			},
			TransactionTypeId: transactionType.Id,
		}

		dataJson, err := util.ToJson(data)
		if err != nil {
			return nil, err
		}

		config := util.NewInlineKeyboardConfig(transactionType.Name, dataJson)
		configs = append(configs, config)
	}

	return util.NewInlineKeyboard(configs, messageContextId, colSize, true), nil
}
