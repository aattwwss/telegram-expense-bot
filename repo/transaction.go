package repo

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/entity"
)

type TransactionRepo struct {
	transactionDao dao.TransactionDAO
}

func NewTransactionRepo(transactionDao dao.TransactionDAO) TransactionRepo {
	return TransactionRepo{transactionDao: transactionDao}
}

func (repo TransactionRepo) Add(ctx context.Context, t domain.Transaction) error {

	err := repo.transactionDao.Insert(ctx, entity.Transaction{
		Id:          t.Id,
		Datetime:    t.Datetime,
		CategoryId:  t.CategoryId,
		Description: t.Description,
		UserId:      t.UserId,
		Amount:      t.Amount.Amount(),
		Currency:    t.Amount.Currency().Code,
	})

	if err != nil {
		return err
	}

	return nil
}
