package util

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
)

type InlineKeyboardConfig struct {
	label string
	data  string
}

func CloseInlineKeyboard(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	editMsgConfig := tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      update.CallbackQuery.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.MessageID,
			ReplyMarkup: nil,
		},
	}
	if _, err := bot.Request(editMsgConfig); err != nil {
		return err
	}
	return nil
}

func NewInlineKeyboardConfig(label string, data string) InlineKeyboardConfig {
	return InlineKeyboardConfig{
		label: label,
		data:  data,
	}
}

func NewInlineKeyboard(configs []InlineKeyboardConfig, colSize int, cancellable bool) [][]tgbotapi.InlineKeyboardButton {
	numOfRows := roundUpDivision(len(configs), colSize)
	var itemsKeyboards [][]tgbotapi.InlineKeyboardButton

	for i := 0; i < numOfRows; i++ {
		row := tgbotapi.NewInlineKeyboardRow()
		for j := 0; j < colSize; j++ {
			itemIndex := colSize*i + j
			if itemIndex == len(configs) {
				break
			}
			config := configs[itemIndex]
			button := tgbotapi.NewInlineKeyboardButtonData(config.label, config.data)
			row = append(row, button)
		}
		itemsKeyboards = append(itemsKeyboards, row)
	}

	if cancellable {
		row := tgbotapi.NewInlineKeyboardRow()
		button := tgbotapi.NewInlineKeyboardButtonData("Cancel", "Cancel||cancel")
		row = append(row, button)
		itemsKeyboards = append(itemsKeyboards, row)
	}

	return itemsKeyboards
}

func roundUpDivision(dividend int, divisor int) int {
	quotient := float64(dividend) / float64(divisor)
	quotientCeiling := math.Ceil(quotient)
	return int(quotientCeiling)
}
