package util

import (
	"math"

	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/enum"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InlineKeyboardConfig struct {
	label string
	data  string
}

func NewPaginationKeyboard(totalCount int, currentOffset int, limit int, messageContextId int, colSize int) ([][]tgbotapi.InlineKeyboardButton, error) {
	var configs []InlineKeyboardConfig

	nextOffset := currentOffset + limit
	if nextOffset < totalCount && limit < totalCount {
		nextButton := domain.PaginationCallback{
			Callback: domain.Callback{
				Type:             enum.Pagination,
				MessageContextId: messageContextId,
			},
			Action: enum.Previous,
			Offset: currentOffset + limit,
			Limit:  limit,
		}

		nextButtonJson, err := ToJson(nextButton)
		if err != nil {
			return nil, err
		}
		configs = append(configs, NewInlineKeyboardConfig(">", nextButtonJson))
	}

	if currentOffset != 0 {
		prevButton := domain.PaginationCallback{
			Callback: domain.Callback{
				Type:             enum.Pagination,
				MessageContextId: messageContextId,
			},
			Action: enum.Next,
			Offset: currentOffset - limit,
			Limit:  limit,
		}

		prevButtonJson, err := ToJson(prevButton)
		if err != nil {
			return nil, err
		}
		configs = append(configs, NewInlineKeyboardConfig("<", prevButtonJson))
	}

	showCancelButton := len(configs) > 0
	return NewInlineKeyboard(configs, messageContextId, colSize, showCancelButton), nil
}

func NewEditEmptyInlineKeyboard(chatId int64, messageId int) tgbotapi.EditMessageReplyMarkupConfig {
	return tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatId,
			MessageID:   messageId,
			ReplyMarkup: nil,
		},
	}
}

func NewInlineKeyboardConfig(label string, data string) InlineKeyboardConfig {
	return InlineKeyboardConfig{
		label: label,
		data:  data,
	}
}

func NewInlineKeyboard(configs []InlineKeyboardConfig, messageContextId int, colSize int, cancellable bool) [][]tgbotapi.InlineKeyboardButton {
	numOfRows := roundUpDivision(len(configs), colSize)
	itemsKeyboards := [][]tgbotapi.InlineKeyboardButton{{}} // important to initiate the inner array to allow empty keyboard

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
		cancelCallback := domain.GenericCallback{
			Callback: domain.Callback{
				Type:             enum.Cancel,
				MessageContextId: messageContextId,
			},
		}
		dataJson, _ := ToJson(cancelCallback)
		row := tgbotapi.NewInlineKeyboardRow()
		button := tgbotapi.NewInlineKeyboardButtonData("Cancel", dataJson)
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
