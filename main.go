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
		case "stat":
			commandHandler.Stat(ctx, &msg, update)
		default:
			commandHandler.Help(ctx, &msg, update)
		}
	} else {
		commandHandler.Transact(ctx, &msg, update)
	}
	botSend(bot, msg)
}

func loadEnv() error {
	appEnv := os.Getenv("APP_ENV")
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

func main() {
	ctx := context.Background()
	err := loadEnv()
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
	categoryDao := dao.NewCategoryDAO(dbLoaded)
	statDao := dao.NewStatDAO(dbLoaded)

	transactionRepo := repo.NewTransactionRepo(transactionDAO)
	userRepo := repo.NewUserRepo(userDAO)
	statRepo := repo.NewStatRepo(statDao)

	commandHandler := handler.NewCommandHandler(userRepo, transactionDAO, categoryDao, statRepo)
	callbackHandler := handler.NewCallbackHandler(userDAO, transactionRepo, categoryDao)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramApiToken)
	if err != nil {
		log.Fatal().Err(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for i := 0; i < 2; i++ {
		go func(bot *tgbotapi.BotAPI, update <-chan tgbotapi.Update) {
			for update := range updates {
				if update.Message != nil {
					handleMessage(ctx, bot, update, &commandHandler)
				} else if update.CallbackQuery != nil {
					handleCallback(ctx, bot, update, &callbackHandler)
				}
			}
		}(bot, updates)
	}
	select {}
}
