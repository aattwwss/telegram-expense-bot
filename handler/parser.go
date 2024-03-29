package handler

import (
	"errors"
	"fmt"
	"regexp"
)

var floatParser = regexp.MustCompile(`^(-?\d[\d]*[.]?[\d{2}]*)`)

// parseFloatStringFromString retrieves a valid float string from a string
func parseFloatStringFromString(s string) (string, error) {
	matches := floatParser.FindAllString(s, -1)
	fmt.Printf("%v", matches)
	if len(matches) == 0 {
		return "", errors.New("no float found in string: " + s)
	}
	return matches[0], nil
}
