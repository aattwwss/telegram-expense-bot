package handler

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
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
	signUpSuccessMsg         = "Welcome to your expense tracker!\nCongratulations! We can get you started right away!\n"
	cannotRecogniseAmountMsg = "I don't recognise that amount of money :(\n"
	descriptionTooLong       = "Sorry, your description (max 20 characters) is too long :( \n"
	transactionListEmptyMsg  = "You have no transactions this month."

	statsHeaderHTMLMsg = "<b>%s %v\n</b>%s\n\n" // E.g. November 2022

	transactionTypeInlineColSize = 2

	listDefaultPageSize   = 10
	exportDefaultPageSize = 1000
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

func (handler CommandHandler) Start(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	teleUser := update.SentFrom()

	dbUser, err := handler.userRepo.FindUserById(ctx, teleUser.ID)
	if err != nil {
		log.Error().Msgf("error finding user: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, errorFindingUserMsg)
		return
	}

	if dbUser != nil {
		log.Info().Msgf("User already exists. id: %v", dbUser.Id)
		util.BotSendMessage(bot, update.Message.Chat.ID, userExistsMsg)
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
		log.Error().Msgf("error adding user: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, errorCreatingUserMsg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, signUpSuccessMsg)
	util.BotSendWrapper(bot, msg)
}

func (handler CommandHandler) Help(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message.HelpMsg)
	util.BotSendWrapper(bot, msg)
}

func (handler CommandHandler) Undo(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userId := update.Message.From.ID
	latestTransaction, err := handler.transactionRepo.FindLastestByUserId(ctx, userId)
	if err != nil {
		log.Error().Msgf("Error finding latest transaction: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	if latestTransaction == nil {
		util.BotSendMessage(bot, update.Message.Chat.ID, message.TransactionLatestNotFound)
	}

	err = handler.transactionRepo.DeleteById(ctx, latestTransaction.Id, userId)
	if err != nil {
		log.Error().Msgf("Error deleting latest transaction: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	text := fmt.Sprintf(message.TransactionDeletedReplyMsg, latestTransaction.Amount.Display(), latestTransaction.Description)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	util.BotSendWrapper(bot, msg)
}

func (handler CommandHandler) StartTransaction(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userId := update.SentFrom().ID
	user, err := handler.userRepo.FindUserById(ctx, userId)
	if err != nil {
		log.Error().Msgf("Error finding user for transact: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}
	if user == nil {
		log.Error().Msgf("User not found for transact: %w", userId)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	floatString, err := util.ParseFloatStringFromString(update.Message.Text)
	if err != nil {
		log.Error().Msgf("%w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, cannotRecogniseAmountMsg)
		return
	}

	stringAfter := util.After(update.Message.Text, floatString)
	if len(strings.TrimSpace(stringAfter)) > 20 {
		util.BotSendMessage(bot, update.Message.Chat.ID, descriptionTooLong)
		return
	}

	contextId, err := handler.messageContextRepo.Add(ctx, update.Message.Chat.ID, update.Message.MessageID, update.Message.Text)
	if err != nil {
		log.Error().Msgf("Add message context error: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	categories, err := handler.categoryRepo.FindAll(ctx)
	if err != nil {
		log.Error().Msgf("FindAll categories error: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	inlineKeyboard, err := newCategoriesKeyboard(categories, contextId, categoriesInlineColSize)
	if err != nil {
		log.Error().Msgf("newCategoriesKeyboard error: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message.TransactionTypeReplyMsg)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}
	util.BotSendWrapper(bot, msg)
}

func (handler CommandHandler) Stats(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userId := update.SentFrom().ID
	user, err := handler.userRepo.FindUserById(ctx, userId)
	if err != nil {
		log.Error().Msgf("Error finding user for stats: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	month, year := util.ParseMonthYearFromMessage(update.Message.Text)

	breakdowns, total, err := handler.transactionRepo.GetTransactionBreakdownByCategory(ctx, month, year, *user)

	if err != nil {
		log.Error().Msgf("Error getting breakdowns: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	header := fmt.Sprintf(statsHeaderHTMLMsg, month.String(), year, total.Display())
	text := header + breakdowns.GetFormattedHTMLMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	util.BotSendWrapper(bot, msg)
}

func (handler CommandHandler) List(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	pageSize := listDefaultPageSize
	userId := update.SentFrom().ID
	user, err := handler.userRepo.FindUserById(ctx, userId)
	if err != nil {
		log.Error().Msgf("Error finding user for stats: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	contextId, err := handler.messageContextRepo.Add(ctx, update.Message.Chat.ID, update.Message.MessageID, update.Message.Text)
	if err != nil {
		log.Error().Msgf("Add message context error: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	month, year := util.ParseMonthYearFromMessage(update.Message.Text)

	transactions, totalCount, err := handler.transactionRepo.ListByMonthAndYear(ctx, month, year, 0, pageSize, *user)
	if err != nil {
		log.Error().Msgf("Error getting list of transactions: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}
	if totalCount == 0 {
		util.BotSendMessage(bot, update.Message.Chat.ID, transactionListEmptyMsg)
		return
	}

	inlineKeyboard, err := util.NewPaginationKeyboard(totalCount, 0, pageSize, contextId, 2)
	if err != nil {
		log.Error().Msgf("Error generating keyboard for transaction pagination: %w", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	text := transactions.GetFormattedHTMLMsg(month, year, user.Location, totalCount, 0, pageSize)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}
	msg.ParseMode = tgbotapi.ModeHTML
	util.BotSendWrapper(bot, msg)
}

func (handler CommandHandler) Export(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	pageSize := exportDefaultPageSize

	userId := update.SentFrom().ID
	user, err := handler.userRepo.FindUserById(ctx, userId)
	if err != nil {
		log.Error().Msgf("Error finding user for stats: %v", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	month, year := util.ParseMonthYearFromMessage(update.Message.Text)
	fileName := fmt.Sprintf("expenses_%02d_%v_*.csv", int(month), year)
	f, err := os.CreateTemp("", fileName)
	defer os.Remove(f.Name())

	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()

	offset := 0
	for {
		transactions, totalCount, err := handler.transactionRepo.ListByMonthAndYear(ctx, month, year, offset, pageSize, *user)
		if offset > totalCount {
			break
		}
		if err != nil {
			log.Error().Msgf("Error finding listing transactions for export: %v", err)
			util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
			return
		}
		for _, t := range transactions {
			data := []string{
				t.Datetime.String(),
				t.CategoryName,
				fmt.Sprintf("%v", t.Amount.AsMajorUnits()),
				t.Amount.Currency().Code,
				t.Description,
			}
			csvWriter.Write(data)
			csvWriter.Flush()
		}
		offset += pageSize
	}

	docMsg := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FilePath(f.Name()))
	docMsg.Caption = fmt.Sprintf("Exported expenses for %s %v", month.String(), year)
	util.BotSendWrapper(bot, docMsg)
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
