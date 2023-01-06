package repo

import (
	"context"
	"time"

	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/entity"
)

type MessageContextRepo struct {
	messageContextDAO dao.MessageContextDAO
}

func NewMessageContextRepo(messageContextDAO dao.MessageContextDAO) MessageContextRepo {
	return MessageContextRepo{messageContextDAO: messageContextDAO}
}

func (repo MessageContextRepo) Add(ctx context.Context, chatId int64, messageId int, message string) (int, error) {

	id, err := repo.messageContextDAO.Insert(ctx, entity.MessageContext{
		ChatId:    chatId,
		MessageId: messageId,
		Message:   message,
		CreatedAt: time.Now(),
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

func (repo MessageContextRepo) DeleteById(ctx context.Context, id int) error {
	err := repo.messageContextDAO.DeleteById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
