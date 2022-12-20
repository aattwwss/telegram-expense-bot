package repo

import (
	"context"

	"github.com/Rhymond/go-money"
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

func (repo TransactionRepo) GetById(ctx context.Context, id int, userId int64) (domain.Transaction, error) {
	e, err := repo.transactionDao.GetById(ctx, id, userId)

	if err != nil {
		return domain.Transaction{}, err
	}

	t := domain.Transaction{
		Id:          e.Id,
		Datetime:    e.Datetime,
		CategoryId:  e.CategoryId,
		Description: e.Description,
		UserId:      e.UserId,
		Amount:      money.New(e.Amount, e.Currency),
	}
	return t, nil
}

func (repo TransactionRepo) FindLastestByUserId(ctx context.Context, userId int64) (*domain.Transaction, error) {
	e, err := repo.transactionDao.FindLatestByUserId(ctx, userId)

	if err != nil {
		return nil, err
	}

	if e == nil {
		return nil, nil
	}

	t := domain.Transaction{
		Id:          e.Id,
		Datetime:    e.Datetime,
		CategoryId:  e.CategoryId,
		Description: e.Description,
		UserId:      e.UserId,
		Amount:      money.New(e.Amount, e.Currency),
	}
	return &t, nil
}

func (repo TransactionRepo) DeleteById(ctx context.Context, id int, userId int64) error {
	err := repo.transactionDao.DeleteById(ctx, id, userId)

	if err != nil {
		return err
	}

	return nil
}
