package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/enum"
	"github.com/aattwwss/telegram-expense-bot/message"
	"github.com/aattwwss/telegram-expense-bot/repo"
	"github.com/aattwwss/telegram-expense-bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type MessageHandler struct {
	transactionRepo     repo.TransactionRepo
	messageContextRepo  repo.MessageContextRepo
	transactionTypeRepo repo.TransactionTypeRepo
	categoryRepo        repo.CategoryRepo
	statRepo            repo.StatRepo
	userRepo            repo.UserRepo
}

func NewMessageHandler(userRepo repo.UserRepo, transactionRepo repo.TransactionRepo, messageContextRepo repo.MessageContextRepo, transactionTypeRepo repo.TransactionTypeRepo, categoryRepo repo.CategoryRepo, statRepo repo.StatRepo) MessageHandler {
	return MessageHandler{
		userRepo:            userRepo,
		transactionRepo:     transactionRepo,
		messageContextRepo:  messageContextRepo,
		transactionTypeRepo: transactionTypeRepo,
		categoryRepo:        categoryRepo,
		statRepo:            statRepo,
	}
}

func (mh MessageHandler) Handle(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	teleUser := update.SentFrom()

	dbUser, err := mh.userRepo.FindUserById(ctx, teleUser.ID)
	if err != nil {
		log.Error().Msgf("error finding user: %v", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, errorFindingUserMsg)
		return
	}
	if dbUser == nil {
		log.Error().Msgf("User not found for transact: %v", teleUser.ID)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	switch dbUser.CurrentContext {
	case enum.Transaction:
		mh.startTransaction(ctx, bot, update)
	case enum.SetTimeZone:
		mh.setTimeZone(ctx, bot, update, *dbUser)
	case enum.SetCurrency:
		mh.setCurrency(ctx, bot, update)
	default:
		mh.startTransaction(ctx, bot, update)
	}
}

func (mh MessageHandler) startTransaction(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	floatString, err := parseFloatStringFromString(update.Message.Text)
	if err != nil {
		log.Error().Msgf("%v", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, cannotRecogniseAmountMsg)
		return
	}

	stringAfter := util.After(update.Message.Text, floatString)
	if len(strings.TrimSpace(stringAfter)) > descLengthLimit {
		util.BotSendMessage(bot, update.Message.Chat.ID, descriptionTooLong)
		return
	}

	contextId, err := mh.messageContextRepo.Add(ctx, update.Message.Chat.ID, update.Message.MessageID, update.Message.Text)
	if err != nil {
		log.Error().Msgf("Add message context error: %v", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	categories, err := mh.categoryRepo.FindAll(ctx)
	if err != nil {
		log.Error().Msgf("FindAll categories error: %v", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	inlineKeyboard, err := newCategoriesKeyboard(categories, contextId, categoriesInlineColSize)
	if err != nil {
		log.Error().Msgf("newCategoriesKeyboard error: %v", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.GenericErrReplyMsg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message.TransactionTypeReplyMsg)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}
	util.BotSendWrapper(bot, msg)
}

func (mh MessageHandler) setTimeZone(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, user domain.User) {
	var tzCallback domain.TimeZoneCallback

	err := json.Unmarshal([]byte(update.CallbackQuery.Data), &tzCallback)
	if err != nil {
		log.Error().Msgf("setTimeZone unmarshall error: %v", err)
		return
	}
	log.Info().Msgf("timezone offset: %v", tzCallback.TzOffset)

	newLoc, err := time.LoadLocation(fmt.Sprintf("Etc/GMT%+d", -1*tzCallback.TzOffset))
	if err != nil {
		log.Error().Msgf("loadLocation error: %v", err)
		util.BotSendMessage(bot, update.Message.Chat.ID, message.InvalidTimeZoneMsg)
		return
	}

	user.Location = newLoc
	mh.userRepo.Update(ctx, user)

	util.BotSendMessage(bot, update.Message.Chat.ID, fmt.Sprintf(message.SetTimeZoneSuccessMsg, offsetToGMT(tzCallback.TzOffset)))
}

func (mh MessageHandler) setCurrency(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
}

func newTimeZoneKeyboard(messageContextId int, colSize int) ([][]tgbotapi.InlineKeyboardButton, error) {
	tzOffsetStart := -12
	tzOffsetEnd := 13

	var configs []util.InlineKeyboardConfig
	for offset := tzOffsetStart; offset <= tzOffsetEnd; offset++ {
		data := domain.TimeZoneCallback{
			Callback: domain.Callback{
				Type:             enum.Category,
				MessageContextId: messageContextId,
			},
			TzOffset: offset,
		}

		dataJson, err := util.ToJson(data)
		if err != nil {
			return nil, err
		}

		config := util.NewInlineKeyboardConfig(offsetToGMT(offset), dataJson)
		configs = append(configs, config)
	}

	return util.NewInlineKeyboard(configs, messageContextId, colSize, true), nil

}

func offsetToGMT(offset int) string {
	return fmt.Sprintf("GMT%+d", offset)
}
