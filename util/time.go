package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type YearMonth struct {
	Month time.Month
	Year  int
}

func (ym YearMonth) String(layout string) (string, error) {
	//2006-01
	if !strings.Contains(layout, "2006") || !strings.Contains(layout, "01") {
		return "", errors.New(fmt.Sprintf("invalid layout %s", layout))
	}
	year := strconv.Itoa(ym.Year)
	layout = strings.ReplaceAll(layout, "2006", year)
	layout = strings.ReplaceAll(layout, "01", fmt.Sprintf("%02d", ym.Month))
	return layout, nil
}
