package util

import (
	"encoding/json"
)

func ToJson[T any](t T) (string, error) {
	s, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(s), nil
}
