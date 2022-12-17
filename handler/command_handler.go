package handler

import (
	"context"
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
	userExistsMsg            = "Welcome back! These are the summary of your transactions: \n"
	errorFindingUserMsg      = "Sorry there is a problem fetching your information.\n"
	errorCreatingUserMsg     = "Sorry there is a problem signing you up.\n"
	signUpSuccessMsg         = "Congratulations! We can get you started right away!\n"
	helpMsg                  = "Type /start to register.\nType /stats to view your last 3 months expenses.\n\nStart recording your expenses by typing the amount you want to save, followed by the description.\n\ne.g. 12.34 Canned pasta"
	cannotRecogniseAmountMsg = "I don't recognise that amount of money :(\n"

	transactionHeaderHTMLMsg = "<b>Summary\n</b>"

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

	user := domain.User{
		Id:       teleUser.ID,
		Locale:   "en",
		Currency: money.GetCurrency(money.SGD),
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

	amountString, err := util.ParseFloatStringFromString(update.Message.Text)
	if err != nil {
		log.Error().Msgf("%v", err)
		msg.Text = cannotRecogniseAmountMsg
		return
	}

	float, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		msg.Text = cannotRecogniseAmountMsg
		return
	}

	amount := money.NewFromFloat(float, user.Currency.Code)
	if err != nil {
		msg.Text = cannotRecogniseAmountMsg
		return
	}

	transactionTypes, err := handler.transactionTypeRepo.GetAll(ctx)
	if err != nil {
		msg.Text = message.GenericErrReplyMsg
		return
	}

	id, err := handler.messageContextRepo.Add(ctx, update.Message.Text)
	if err != nil {
		msg.Text = message.GenericErrReplyMsg
		return
	}

	inlineKeyboard, err := newTransactionTypesKeyboard(transactionTypes, id, transactionTypeInlineColSize)
	if err != nil {
		msg.Text = message.GenericErrReplyMsg
		return
	}

	msg.Text = fmt.Sprintf(message.TransactionTypeReplyMsg+"%v", amount.AsMajorUnits())
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
	if user == nil {
		log.Error().Msgf("User not found for stats: %v", userId)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	now := time.Now()
	monthTo := util.YearMonth{
		Month: now.Month(),
		Year:  now.Year(),
	}

	from := now.AddDate(0, -2, 0)
	monthFrom := util.YearMonth{
		Month: from.Month(),
		Year:  from.Year(),
	}

	param := repo.GetMonthlySearchParam{
		Location:  *user.Location,
		MonthFrom: monthFrom,
		MonthTo:   monthTo,
		UserId:    userId,
	}

	summaries, err := handler.statRepo.GetMonthly(ctx, param)
	if err != nil {
		log.Error().Msgf("%v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}

	msg.Text = transactionHeaderHTMLMsg
	msg.Text += summaries.GenerateReportText()

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
