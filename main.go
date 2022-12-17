package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/aattwwss/telegram-expense-bot/config"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/db"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/handler"
	"github.com/aattwwss/telegram-expense-bot/message"
	"github.com/aattwwss/telegram-expense-bot/repo"
	"github.com/aattwwss/telegram-expense-bot/util"
	"github.com/caarlos0/env/v6"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func botSend(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		log.Error().Msgf("handleCallback error: %v", err)
	}
}

func getCallbackType(callbackData string) (string, error) {
	var genericCallback domain.GenericCallback
	err := json.Unmarshal([]byte(callbackData), &genericCallback)
	if err != nil {
		return "", err
	}
	return genericCallback.Type, nil
}

func handleCallback(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, callbackHandler *handler.CallbackHandler) {
	err := util.CloseInlineKeyboard(bot, update)
	if err != nil {
		log.Error().Msgf("handleCallback error: %v", err)
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
	callbackType, err := getCallbackType(update.CallbackQuery.Data)
	if err != nil {
		msg.Text = message.GenericErrReplyMsg
		botSend(bot, msg)
	}

	switch callbackType {
	case "TransactionType":
		callbackHandler.FromTransactionType(ctx, &msg, update.CallbackQuery)
	case "Category":
		callbackHandler.FromCategory(ctx, &msg, update.CallbackQuery)
	case "Cancel":
		callbackHandler.FromCancel(ctx, &msg, update.CallbackQuery)
		return
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
		commandHandler.StartTransaction(ctx, &msg, update)
	}
	botSend(bot, msg)
}

func loadEnv() error {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" || appEnv == "DEV" {
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

func runWebhook(bot *tgbotapi.BotAPI, cfg config.EnvConfig) tgbotapi.UpdatesChannel {
	log.Info().Msg("Running on webhook!")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })

	go http.ListenAndServe(cfg.AppHost+":"+cfg.AppPort, nil)
	time.Sleep(200 * time.Millisecond)

	webhook, err := tgbotapi.NewWebhook(cfg.WebhookHost + "/" + bot.Token)
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

	s, _ := json.Marshal(info)
	log.Info().Msgf("Telegram callback info: %s", s)

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
	transactionDao := dao.NewTransactionDao(dbLoaded)
	messageContextDao := dao.NewMessageContextDao(dbLoaded)
	transactionTypeDao := dao.NewTransactionTypeDAO(dbLoaded)
	categoryDao := dao.NewCategoryDAO(dbLoaded)
	statDao := dao.NewStatDAO(dbLoaded)

	transactionRepo := repo.NewTransactionRepo(transactionDao)
	messageContextRepo := repo.NewMessageContextRepo(messageContextDao)
	transactionTypeRepo := repo.NewTransactionTypeRepo(transactionTypeDao)
	userRepo := repo.NewUserRepo(userDAO)
	statRepo := repo.NewStatRepo(statDao)
	categoryRepo := repo.NewCategoryRepo(categoryDao)

	commandHandler := handler.NewCommandHandler(userRepo, transactionRepo, messageContextRepo, transactionTypeRepo, categoryRepo, statRepo)
	callbackHandler := handler.NewCallbackHandler(userRepo, transactionRepo, messageContextRepo, transactionTypeRepo, categoryRepo)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramApiToken)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	var updates tgbotapi.UpdatesChannel

	if cfg.WebhookEnabled {
		updates = runWebhook(bot, cfg)
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
