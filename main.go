package main

import (
	"context"
	"errors"
	"github.com/aattwwss/telegram-expense-bot/config"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/db"
	"github.com/aattwwss/telegram-expense-bot/handler"
	"github.com/aattwwss/telegram-expense-bot/message"
	"github.com/aattwwss/telegram-expense-bot/repo"
	"github.com/caarlos0/env/v6"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

func botSend(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		log.Error().Msgf("handleCallback error: %v", err)
	}
}

func decodeCallbackData(update tgbotapi.Update) (string, string, error) {
	dataArr := strings.Split(update.CallbackQuery.Data, "||")
	if len(dataArr) != 2 {
		return "", "", errors.New("decodeCallbackData error")
	}
	return dataArr[0], dataArr[1], nil
}

func closeInlineKeyboard(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	editMsgConfig := tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      update.CallbackQuery.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.MessageID,
			ReplyMarkup: nil,
		},
	}
	if _, err := bot.Request(editMsgConfig); err != nil {
		log.Error().Msgf("handleMessage error: %v", err)
	}
}

func handleCallback(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, callbackHandler *handler.CallbackHandler) {
	closeInlineKeyboard(bot, update)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
	typeName, data, err := decodeCallbackData(update)
	if err != nil {
		msg.Text = message.GenericErrReplyMsg
		botSend(bot, msg)
	}

	switch typeName {
	case "TransactionType":
		callbackHandler.FromTransactionType(ctx, &msg, update.CallbackQuery, data)
	case "Category":
		callbackHandler.FromCategory(ctx, &msg, update.CallbackQuery, data)
	default:
		log.Error().Msg("handleCallback error: unrecognised callback")
		msg.Text = message.GenericErrReplyMsg
	}

	botSend(bot, msg)
}
func handleMessage(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, commandHandler *handler.CommandHandler) {
	log.Info().Msgf("Received: %v", update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			commandHandler.Start(ctx, &msg, update)
		case "help":
			commandHandler.Help(ctx, &msg, update)
		case "stats":
			commandHandler.Stats(ctx, &msg, update)
		default:
			commandHandler.Help(ctx, &msg, update)
		}
	} else {
		commandHandler.Transact(ctx, &msg, update)
	}
	botSend(bot, msg)
}

func loadEnv(appEnv string) error {
	if appEnv == "" || appEnv == "dev" {
		err := godotenv.Load(".env.local")
		if err != nil {
			return err
		}
	}

	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}

func runWebhook(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	log.Info().Msg("Running on webhook!")
	webhook, err := tgbotapi.NewWebhook("http://localhost:8123/" + bot.Token)
	if err != nil {
		log.Fatal().Err(err)
	}

	_, err = bot.Request(webhook)
	if err != nil {
		log.Fatal().Err(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("Telegram callback failed: %s", info)

	if info.LastErrorDate != 0 {
		log.Info().Msgf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return bot.ListenForWebhook("/" + bot.Token)
}

func runPolling(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	log.Info().Msg("Running on polling!")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return bot.GetUpdatesChan(u)
}

func main() {
	ctx := context.Background()
	appEnv := os.Getenv("APP_ENV")
	err := loadEnv(appEnv)
	if err != nil {
		log.Fatal().Err(err)
	}
	cfg := config.EnvConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err)
	}
	dbLoaded, _ := db.LoadDB(ctx, cfg)

	userDAO := dao.NewUserDao(dbLoaded)
	transactionDAO := dao.NewTransactionDAO(dbLoaded)
	transactionTypeDAO := dao.NewTransactionTypeDAO(dbLoaded)
	categoryDao := dao.NewCategoryDAO(dbLoaded)
	statDao := dao.NewStatDAO(dbLoaded)

	transactionRepo := repo.NewTransactionRepo(transactionDAO)
	transactionTypeRepo := repo.NewTransactionTypeRepo(transactionTypeDAO)
	userRepo := repo.NewUserRepo(userDAO)
	statRepo := repo.NewStatRepo(statDao)
	categoryRepo := repo.NewCategoryRepo(categoryDao)

	commandHandler := handler.NewCommandHandler(userRepo, transactionRepo, transactionTypeRepo, categoryRepo, statRepo)
	callbackHandler := handler.NewCallbackHandler(userRepo, transactionRepo, transactionTypeRepo, categoryRepo)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramApiToken)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	var updates tgbotapi.UpdatesChannel

	if strings.EqualFold(appEnv, "prod") {
		updates = runWebhook(bot)
	} else {
		updates = runPolling(bot)
	}

	processUpdate := func(bot *tgbotapi.BotAPI, update <-chan tgbotapi.Update) {
		for update := range updates {
			if update.Message != nil {
				handleMessage(ctx, bot, update, &commandHandler)
			} else if update.CallbackQuery != nil {
				handleCallback(ctx, bot, update, &callbackHandler)
			}
		}
	}

	for i := 0; i < cfg.NumRoutines; i++ {
		go processUpdate(bot, updates)
	}
	select {}

}
