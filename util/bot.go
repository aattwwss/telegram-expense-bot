package util

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func BotSendWrapper(bot *tgbotapi.BotAPI, chattables ...tgbotapi.Chattable) {
	for _, c := range chattables {
		_, err := bot.Send(c)
		if err != nil {
			log.Error().Msgf("bot send chattable error: %w", err)
			return
		}
	}
}

func BotSendMessage(bot *tgbotapi.BotAPI, chatId int64, message string) {
	m := tgbotapi.NewMessage(chatId, message)
	BotSendWrapper(bot, m)
}
