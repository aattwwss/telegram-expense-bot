package entity

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"time"
)

type User struct {
	Id        int64
	IsBot     bool
	FirstName string
	LastName  *string
	Username  *string
}

func (b1 User) Serialize() (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(b1)
	if err != nil {
		return "", errors.New("user serialization failed")
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

type Transaction struct {
	Id          int64
	Datetime    time.Time
	CategoryId  int
	Description string
	UserId      int64
	Amount      int64
	Currency    string
}

type Category struct {
	Id                int64
	Name              string
	TransactionTypeId int64
}
