package handler

import (
	"context"
	"fmt"
	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/entity"
	"github.com/aattwwss/telegram-expense-bot/message"
	"github.com/aattwwss/telegram-expense-bot/repo"
	"github.com/aattwwss/telegram-expense-bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	userExistsMsg            = "Welcome back! These are the summary of your transactions: \n"
	errorFindingUserMsg      = "Sorry there is a problem fetching your information.\n"
	errorCreatingUserMsg     = "Sorry there is a problem signing you up.\n"
	signUpSuccessMsg         = "Congratulations! We can get you started right away!\n"
	helpMsg                  = "Type /start to register.\nType <category>, <price>, [date]\n"
	cannotRecogniseAmountMsg = "I don't recognise that amount of money :(\n"

	transactionHeaderHTMLMsg  = "<b>Summary\n</b>"
	monthYearHeaderHTMLMsg    = "<code>\n%v %v\n</code>"
	transactionSummaryHTMLMsg = "<code>%v:%s %v\n</code>"
	transactionTotalHTMLMsg   = "<code>🟡 Total: %v\n</code>"

	categoriesInlineColNum = 3
)

type CommandHandler struct {
	//userDao        dao.UserDAO
	transactionDao dao.TransactionDAO
	categoryDao    dao.CategoryDAO
	statDao        dao.StatDAO
	userRepo       repo.UserRepo
}

func NewCommandHandler(userRepo repo.UserRepo, transactionDao dao.TransactionDAO, categoryDao dao.CategoryDAO, statDao dao.StatDAO) CommandHandler {
	return CommandHandler{
		userRepo:       userRepo,
		transactionDao: transactionDao,
		categoryDao:    categoryDao,
		statDao:        statDao,
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

	err = handler.userRepo.Add	a(ctx, user)
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
	param := dao.GetMonthlySearchParam{
		Location:           *loc,
		PrevMonthIntervals: 3,
		UserId:             userId,
	}
	summaries, err := handler.statDao.GetMonthly(ctx, param)
	if err != nil {
		log.Error().Msgf("%v", err)
		msg.Text = message.GenericErrReplyMsg
		return
	}
	formatTransactionLabelInSummaries(summaries)
	msg.Text = transactionHeaderHTMLMsg

	currMonth := ""
	var totalAmountForTheMonth int64

	for i, summary := range summaries {
		month := summary.Month.String()[:3]
		if currMonth != month {
			msg.Text += fmt.Sprintf(monthYearHeaderHTMLMsg, month, summary.Year)
			currMonth = month
			totalAmountForTheMonth = 0
		}

		totalAmountForTheMonth += summary.Amount * summary.Multiplier
		moneyAmount := money.New(summary.Amount, money.SGD)
		msg.Text += fmt.Sprintf(transactionSummaryHTMLMsg, summary.TransactionTypeLabel, strings.Repeat(" ", summary.GetSpacesToPad()), moneyAmount.Display())

		if i == len(summaries)-1 || summaries[i+1].Month.String()[:3] != currMonth {
			msg.Text += fmt.Sprintf(transactionTotalHTMLMsg, money.New(totalAmountForTheMonth, money.SGD).Display())
		}
	}

	msg.ParseMode = tgbotapi.ModeHTML
	return
}

func formatTransactionLabelInSummaries(summaries []*entity.MonthlySummary) {
	longestLabel := 0
	for _, summary := range summaries {
		lengthOfLabel := len(summary.TransactionTypeLabel)
		if lengthOfLabel > longestLabel {
			longestLabel = lengthOfLabel
		}
	}

	for i := range summaries {
		label := summaries[i].TransactionTypeLabel
		summaries[i].SetSpacesToPad(longestLabel - len(label))
	}
}
