package repo

import (
	"context"

	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/entity"
)

type TransactionTypeRepo struct {
	transactionTypeDao dao.TransactionTypeDAO
}

func NewTransactionTypeRepo(transactionTypeDao dao.TransactionTypeDAO) TransactionTypeRepo {
	return TransactionTypeRepo{transactionTypeDao: transactionTypeDao}
}

func (repo TransactionTypeRepo) GetAll(ctx context.Context) ([]*entity.TransactionType, error) {
	return repo.transactionTypeDao.GetAll(ctx)
}

func (repo TransactionTypeRepo) GetById(ctx context.Context, id int) (*entity.TransactionType, error) {
	e, err := repo.transactionTypeDao.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return e, nil
}
