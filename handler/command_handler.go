package handler

import (
	"context"
	"fmt"
	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/message"
	"github.com/aattwwss/telegram-expense-bot/repo"
	"github.com/aattwwss/telegram-expense-bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"math"
	"strconv"
	"time"
)

const (
	userExistsMsg            = "Welcome back! These are the summary of your transactions: \n"
	errorFindingUserMsg      = "Sorry there is a problem fetching your information.\n"
	errorCreatingUserMsg     = "Sorry there is a problem signing you up.\n"
	signUpSuccessMsg         = "Congratulations! We can get you started right away!\n"
	helpMsg                  = "Type /start to register.\nType <category>, <price>, [date]\n"
	cannotRecogniseAmountMsg = "I don't recognise that amount of money :(\n"

	transactionHeaderHTMLMsg = "<b>Summary\n</b>"

	categoriesInlineColNum = 3
)

type CommandHandler struct {
	transactionRepo repo.TransactionRepo
	categoryDao     dao.CategoryDAO

	statRepo repo.StatRepo
	userRepo repo.UserRepo
}

func NewCommandHandler(userRepo repo.UserRepo, transactionRepo repo.TransactionRepo, categoryDao dao.CategoryDAO, statRepo repo.StatRepo) CommandHandler {
	return CommandHandler{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		categoryDao:     categoryDao,
		statRepo:        statRepo,
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

	loc, _ := time.LoadLocation("Asia/Singapore")

	user := domain.User{
		Id:       teleUser.ID,
		Locale:   "en",
		Currency: money.GetCurrency("SGD"),
		Location: loc,
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

func (handler CommandHandler) Transact(ctx context.Context, msg *tgbotapi.MessageConfig, update tgbotapi.Update) {
	float, err := strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		msg.Text = cannotRecogniseAmountMsg
		return
	}

	amount := money.NewFromFloat(float, money.SGD)
	msg.Text = fmt.Sprintf(message.TransactionReplyMsg+"%v", amount.AsMajorUnits())

	categories, err := handler.categoryDao.FindByTransactionTypeId(ctx, 1)
	if err != nil {
		log.Error().Err(err)
		msg.Text = "Sorry we cannot handle your transaction right now :("
		return
	}

	numOfRows := roundUpDivision(len(categories), categoriesInlineColNum)
	var categoriesKeyboards [][]tgbotapi.InlineKeyboardButton

	for i := 0; i < numOfRows; i++ {
		row := tgbotapi.NewInlineKeyboardRow()
		for j := 0; j < categoriesInlineColNum; j++ {
			catIndex := categoriesInlineColNum*i + j
			if catIndex == len(categories) {
				break
			}
			category := categories[catIndex]
			serializedCategory := util.CallbackDataSerialize(*category, category.Id)
			button := tgbotapi.NewInlineKeyboardButtonData(category.Name, serializedCategory)
			row = append(row, button)
		}
		categoriesKeyboards = append(categoriesKeyboards, row)
	}

	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: categoriesKeyboards}
	return
}

func roundUpDivision(dividend int, divisor int) int {
	quotient := float64(dividend) / float64(divisor)
	quotientCeiling := math.Ceil(quotient)
	return int(quotientCeiling)
}

func (handler CommandHandler) Stat(ctx context.Context, msg *tgbotapi.MessageConfig, update tgbotapi.Update) {
	userId := update.SentFrom().ID
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil || loc == nil {
		msg.Text = message.GenericErrReplyMsg
		return
	}
	param := repo.GetMonthlySearchParam{
		Location:           *loc,
		PrevMonthIntervals: 3,
		UserId:             userId,
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
