package util

import (
	"errors"
	"regexp"
)

var floatParser = regexp.MustCompile(`-?\d[\d,]*[.]?[\d{2}]*`)

func ParseFloatStringFromString(s string) (string, error) {
	matches := floatParser.FindAllString(s, -1)
	if len(matches) == 0 {
		return "", errors.New("no float found in string: " + s)
	}
	return matches[0], nil
}
