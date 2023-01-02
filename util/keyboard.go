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

	if currentOffset < totalCount && limit < totalCount {
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
		// symbol is reversed because we are reverse sorting the item in transaction list in reversed
		configs = append(configs, NewInlineKeyboardConfig("<", nextButtonJson))
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
		configs = append(configs, NewInlineKeyboardConfig(">", prevButtonJson))
	}
	return NewInlineKeyboard(configs, messageContextId, colSize, true), nil
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

func NewInlineKeyboard(configs []InlineKeyboardConfig, messageContextId int, colSize int, cancellable bool) [][]tgbotapi.InlineKeyboardButton {
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
