package util

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
)

type InlineKeyboardConfig struct {
	label string
	data  string
}

func NewInlineKeyboardConfig(label string, data string) InlineKeyboardConfig {
	return InlineKeyboardConfig{
		label: label,
		data:  data,
	}
}

func NewInlineKeyboard(configs []InlineKeyboardConfig, colSize int) [][]tgbotapi.InlineKeyboardButton {
	numOfRows := roundUpDivision(len(configs), colSize)
	var categoriesKeyboards [][]tgbotapi.InlineKeyboardButton

	for i := 0; i < numOfRows; i++ {
		row := tgbotapi.NewInlineKeyboardRow()
		for j := 0; j < colSize; j++ {
			catIndex := colSize*i + j
			if catIndex == len(configs) {
				break
			}
			config := configs[catIndex]
			button := tgbotapi.NewInlineKeyboardButtonData(config.label, config.data)
			row = append(row, button)
		}
		categoriesKeyboards = append(categoriesKeyboards, row)
	}
	return categoriesKeyboards
}

func roundUpDivision(dividend int, divisor int) int {
	quotient := float64(dividend) / float64(divisor)
	quotientCeiling := math.Ceil(quotient)
	return int(quotientCeiling)
}
