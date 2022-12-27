package handler

import (
	"context"
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
	userExistsMsg            = "Welcome back! These are the summary of your transactions: \n"
	errorFindingUserMsg      = "Sorry there is a problem fetching your information.\n"
	errorCreatingUserMsg     = "Sorry there is a problem signing you up.\n"
	signUpSuccessMsg         = "Congratulations! We can get you started right away!\n"
	helpMsg                  = "Type /start to register.\nType /stats to view your last 3 months expenses.\n\nStart recording your expenses by typing the amount you want to save, followed by the description.\n\ne.g. 12.34 Canned pasta"
	cannotRecogniseAmountMsg = "I don't recognise that amount of money :(\n"
	descriptionTooLong       = "Sorry, your description is too long :(\n"

	transactionHeaderHTMLMsg = "<b>Summary\n\n</b>"

	transactionTypeInlineColSize = 2
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
	msg.Text = helpMsg
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
	if len(strings.TrimSpace(stringAfter)) > 255 {
		msg.Text = descriptionTooLong
		return
	}

	contextId, err := handler.messageContextRepo.Add(ctx, update.Message.Text)
	if err != nil {
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

	month, year := parseMonthYearFromStatsMessage(update.Message.Text)

	breakdowns, err := handler.transactionRepo.GetTransactionBreakdownByCategory(ctx, month, year, *user)

	if err != nil {
		log.Error().Msgf("Error getting breakdowns: %v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	log.Info().Msgf("breakdowns: %v", breakdowns)
	msg.Text = transactionHeaderHTMLMsg + breakdowns.GetFormattedHTMLText()
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

func parseMonthYearFromStatsMessage(s string) (time.Month, int) {
	now := time.Now()
	month := now.Month()
	year := now.Year()
	arr := strings.Split(s, " ")
	if len(arr) == 2 {
		return parseMonthFromString(arr[1]), year
	}
	if len(arr) == 3 {
		y, err := strconv.Atoi(arr[2])
		if err != nil {
			y = year
		}
		return parseMonthFromString(arr[1]), y
	}
	return month, year
}

func parseMonthFromString(s string) time.Month {
	if s == "1" || strings.EqualFold(s, "jan") || strings.EqualFold(s, "january") {
		return time.January
	}
	if s == "2" || strings.EqualFold(s, "feb") || strings.EqualFold(s, "february") {
		return time.February
	}
	if s == "3" || strings.EqualFold(s, "mar") || strings.EqualFold(s, "march") {
		return time.March
	}
	if s == "4" || strings.EqualFold(s, "apr") || strings.EqualFold(s, "april") {
		return time.April
	}
	if s == "5" || strings.EqualFold(s, "may") || strings.EqualFold(s, "may") {
		return time.May
	}
	if s == "6" || strings.EqualFold(s, "jun") || strings.EqualFold(s, "june") {
		return time.June
	}
	if s == "7" || strings.EqualFold(s, "jul") || strings.EqualFold(s, "july") {
		return time.July
	}
	if s == "8" || strings.EqualFold(s, "aug") || strings.EqualFold(s, "august") {
		return time.August
	}
	if s == "9" || strings.EqualFold(s, "sep") || strings.EqualFold(s, "september") {
		return time.September
	}
	if s == "10" || strings.EqualFold(s, "oct") || strings.EqualFold(s, "october") {
		return time.October
	}
	if s == "11" || strings.EqualFold(s, "nov") || strings.EqualFold(s, "november") {
		return time.November
	}
	if s == "12" || strings.EqualFold(s, "dec") || strings.EqualFold(s, "december") {
		return time.December
	}
	return time.Now().Month()
}
