package main

import (
	"github.com/aattwwss/telegram-expense-bot/config"
	"github.com/caarlos0/env/v6"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
)

func handleFunc(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil { // ignore any non-Message updates
		return
	}

	if !update.Message.IsCommand() { // ignore any non-command Messages
		return
	}

	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	// Extract the command from the Message.
	switch update.Message.Command() {
	case "start":
		msg.Text = "Welcome to your expense tracker!"
	case "help":
		msg.Text = "I understand /sayhi and /status."
	case "sayhi":
		msg.Text = "Hi :)"
	case "status":
		msg.Text = "I'm ok."
	default:
		msg.Text = update.Message.Command()
	}

	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func main() {
	//ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}
	//dbLoaded, _ := db.LoadDB(ctx, cfg)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramApiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for i := 0; i < 5; i++ {
		go func(bot *tgbotapi.BotAPI, update <-chan tgbotapi.Update) {
			for update := range updates {
				handleFunc(bot, update)
			}
		}(bot, updates)
	}
	select {}
}
