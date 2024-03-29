package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"

	"github.com/aattwwss/telegram-expense-bot/config"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/db"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/enum"
	"github.com/aattwwss/telegram-expense-bot/handler"
	"github.com/aattwwss/telegram-expense-bot/message"
	"github.com/aattwwss/telegram-expense-bot/repo"
	"github.com/aattwwss/telegram-expense-bot/util"
	"github.com/caarlos0/env/v6"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func getCallbackType(callbackData string) (enum.CallbackType, error) {
	var genericCallback domain.GenericCallback
	err := json.Unmarshal([]byte(callbackData), &genericCallback)
	if err != nil {
		return "", err
	}
	return genericCallback.Type, nil
}

func handleCallback(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, callbackHandler *handler.CallbackHandler) {
	go util.BotDeleteMessage(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)

	callbackType, err := getCallbackType(update.CallbackQuery.Data)
	if err != nil {
		log.Error().Msg("handleCallback getCallbackType error: unrecognised callback")
		util.BotSendMessage(bot, update.CallbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
	}

	switch callbackType {
	case enum.Category:
		callbackHandler.FromCategory(ctx, bot, update.CallbackQuery)
	case enum.Pagination:
		callbackHandler.FromPagination(ctx, bot, update.CallbackQuery)
	case enum.Undo:
		callbackHandler.FromUndo(ctx, bot, update.CallbackQuery)
	case enum.Cancel:
		callbackHandler.FromCancel(ctx, bot, update.CallbackQuery)
	default:
		log.Error().Msg("handleCallback error: unrecognised callback")
		util.BotSendMessage(bot, update.CallbackQuery.Message.Chat.ID, message.GenericErrReplyMsg)
	}
}

func handleMessage(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, commandHandler *handler.CommandHandler) {
	log.Info().Msgf("Received: %v", update.Message.Text)

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			commandHandler.Start(ctx, bot, update)
		case "help":
			commandHandler.Help(ctx, bot, update)
		case "stats":
			commandHandler.Stats(ctx, bot, update)
		case "undo":
			commandHandler.Undo(ctx, bot, update)
		case "list":
			commandHandler.List(ctx, bot, update)
		case "export":
			commandHandler.Export(ctx, bot, update)
		default:
			commandHandler.Help(ctx, bot, update)
		}
	} else {
		commandHandler.StartTransaction(ctx, bot, update)
	}
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
	lastErrorTime := time.UnixMilli(int64(info.LastErrorDate * 1000))
	if info.LastErrorDate != 0 {
		log.Info().Msgf("Last error date: %s Telegram callback failed: %s", lastErrorTime.Format(time.RFC3339), info.LastErrorMessage)
	}

	return bot.ListenForWebhook("/" + bot.Token)
}

func runPolling(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	log.Info().Msg("Running on polling!")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return bot.GetUpdatesChan(u)
}

type telegramHook struct {
	token           string
	chatId          string
	applicationName string
}

func (h telegramHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.ErrorLevel || level == zerolog.FatalLevel || level == zerolog.PanicLevel {
		apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", h.token)

		logMsg := fmt.Sprintf("<b>%s</b><code>\n\nlevel: %s\ntime : %s\nmsg  : %s</code>", h.applicationName, level, time.Now().Format("2006-01-02 15:04:05.000"), msg)
		requestBody := fmt.Sprintf(`{"chat_id": "%s", "text": "%s", "parse_mode": "HTML"}`, h.chatId, logMsg)
		resp, _ := http.Post(apiURL, "application/json", bytes.NewBuffer([]byte(requestBody)))
		defer resp.Body.Close()
	}
}

func newTelegramHook(token, chatId string) zerolog.Hook {
	telegramHook := telegramHook{
		token:           token,
		chatId:          chatId,
		applicationName: "expense-tracker-bot",
	}
	return telegramHook
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
	if cfg.LogTelegramToken != "" || cfg.LogTelegramChatId != "" {
		telegramHook := newTelegramHook(cfg.LogTelegramToken, cfg.LogTelegramChatId)
		log.Logger = log.Hook(telegramHook)
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
