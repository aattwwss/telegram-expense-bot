package handler

import (
	"context"
	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/entity"
	"github.com/aattwwss/telegram-expense-bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

const (
	userExistsMsg        = "Welcome back! These are the summary of your transactions: \n"
	errorFindingUserMsg  = "Sorry there is a problem fetching your information.\n"
	errorCreatingUserMsg = "Sorry there is a problem signing you up.\n"
	signUpSuccessMsg     = "Congratulations! We can get you started right away!\n"
	registerHereMsg      = "Looks like you have not registered in our system. Type /start to register!\n"
	helpMsg              = "Type /start to register.\nType <category>, <price>, [date]\n"
)

type CommandHandler struct {
	userDao        dao.UserDAO
	transactionDao dao.TransactionDAO
}

func NewCommandHandler(userDao dao.UserDAO, transactionDao dao.TransactionDAO) CommandHandler {
	return CommandHandler{
		userDao:        userDao,
		transactionDao: transactionDao,
	}
}

func (handler CommandHandler) Start(ctx context.Context, msg *tgbotapi.MessageConfig, update tgbotapi.Update) {
	msg.Text = "Welcome to your expense tracker!\n"

	teleUser := update.SentFrom()

	dbUser, err := handler.userDao.FindUserById(ctx, teleUser.ID)
	if err != nil {
		msg.Text += errorFindingUserMsg
		return
	}

	if dbUser != nil {
		msg.Text += userExistsMsg
		return
	}

	entityUser := entity.User{
		Id:        teleUser.ID,
		IsBot:     teleUser.IsBot,
		FirstName: teleUser.FirstName,
		LastName:  util.Ptr(teleUser.LastName),
		Username:  util.Ptr(teleUser.UserName),
	}
	err = handler.userDao.Insert(ctx, entityUser)
	if err != nil {
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

//var numericKeyboard = tgbotapi.NewOneTimeReplyKeyboard(
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton("1"),
//		tgbotapi.NewKeyboardButton("2"),
//		tgbotapi.NewKeyboardButton("3"),
//	),
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton("4"),
//		tgbotapi.NewKeyboardButton("5"),
//		tgbotapi.NewKeyboardButton("6"),
//	),
//)

var inlineNumericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

func (handler CommandHandler) Transact(ctx context.Context, msg *tgbotapi.MessageConfig, update tgbotapi.Update) {
	float, err := strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		msg.Text = "not correct money format :("
		return
	}

	amount := money.NewFromFloat(float, money.SGD)
	msg.Text = "Select the categories this amount belongs to." + amount.Display()
	msg.ReplyMarkup = inlineNumericKeyboard
	return
}
