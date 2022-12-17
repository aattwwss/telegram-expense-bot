package repo

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/entity"
)

type MessageContextRepo struct {
	messageContextDAO dao.MessageContextDAO
}

func NewMessageContextRepo(messageContextDAO dao.MessageContextDAO) MessageContextRepo {
	return MessageContextRepo{messageContextDAO: messageContextDAO}
}

func (repo MessageContextRepo) Add(ctx context.Context, message string) (int, error) {

	id, err := repo.messageContextDAO.Insert(ctx, entity.MessageContext{
		Message: message,
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo MessageContextRepo) GetMessageById(ctx context.Context, id int) (string, error) {
	e, err := repo.messageContextDAO.GetById(ctx, id)
	if err != nil {
		return "", err
	}
	return e.Message, nil
}
