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

// parseMonthYearFromMessage returns the month and year representation from the string,
// any error returns the current month or year
func ParseMonthYearFromMessage(s string) (time.Month, int) {
	now := time.Now()
	month := now.Month()
	year := now.Year()
	arr := strings.Split(s, " ")
	if len(arr) == 2 {
		return ParseMonthFromString(arr[1]), year
	}
	if len(arr) == 3 {
		y, err := strconv.Atoi(arr[2])
		if err != nil {
			y = year
		}
		return ParseMonthFromString(arr[1]), y
	}
	return month, year
}

// ParseMonthFromString trys to return the month given a string, else it returns the current month.
func ParseMonthFromString(s string) time.Month {
	if s == "1" || strings.EqualFold(s, "jan") || strings.EqualFold(s, "january") {
		return time.January
	}
	if s == "2" || strings.EqualFold(s, "feb") || strings.EqualFold(s, "february") {
		return time.February
	}
	if s == "3" || strings.EqualFold(s, "mar") || strings.EqualFold(s, "march") {
		return time.March
	}
	if s == "4" || strings.EqualFold(s, "apr") || strings.EqualFold(s, "april") {
		return time.April
	}
	if s == "5" || strings.EqualFold(s, "may") || strings.EqualFold(s, "may") {
		return time.May
	}
	if s == "6" || strings.EqualFold(s, "jun") || strings.EqualFold(s, "june") {
		return time.June
	}
	if s == "7" || strings.EqualFold(s, "jul") || strings.EqualFold(s, "july") {
		return time.July
	}
	if s == "8" || strings.EqualFold(s, "aug") || strings.EqualFold(s, "august") {
		return time.August
	}
	if s == "9" || strings.EqualFold(s, "sep") || strings.EqualFold(s, "september") {
		return time.September
	}
	if s == "10" || strings.EqualFold(s, "oct") || strings.EqualFold(s, "october") {
		return time.October
	}
	if s == "11" || strings.EqualFold(s, "nov") || strings.EqualFold(s, "november") {
		return time.November
	}
	if s == "12" || strings.EqualFold(s, "dec") || strings.EqualFold(s, "december") {
		return time.December
	}
	return time.Now().Month()
}
