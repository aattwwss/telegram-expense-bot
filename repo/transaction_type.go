package repo

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/domain"
)

type TransactionTypeRepo struct {
	transactionTypeDao dao.TransactionTypeDAO
}

func NewTransactionTypeRepo(transactionTypeDao dao.TransactionTypeDAO) TransactionTypeRepo {
	return TransactionTypeRepo{transactionTypeDao: transactionTypeDao}
}

func (repo TransactionTypeRepo) GetAll(ctx context.Context) ([]domain.TransactionType, error) {
	entities, err := repo.transactionTypeDao.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var types []domain.TransactionType

	for _, e := range entities {
		transactionType := domain.TransactionType{
			Id:         e.Id,
			Name:       e.Name,
			Multiplier: e.Multiplier,
			ReplyText:  e.ReplyText,
		}
		types = append(types, transactionType)
	}

	return types, nil
}

func (repo TransactionTypeRepo) GetById(ctx context.Context, id int) (domain.TransactionType, error) {
	e, err := repo.transactionTypeDao.GetById(ctx, id)
	if err != nil {
		return domain.TransactionType{}, err
	}

	transactionType := domain.TransactionType{
		Id:         e.Id,
		Name:       e.Name,
		Multiplier: e.Multiplier,
		ReplyText:  e.ReplyText,
	}

	return transactionType, nil
}
