package repo

import (
	"context"
	"fmt"
	"time"

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

func (repo TransactionRepo) GetTransactionBreakdownByCategory(ctx context.Context, month time.Month, year int, user domain.User) (domain.Breakdowns, error) {
	breakdowns := domain.Breakdowns{}

	dateFromString := fmt.Sprintf("%v-%02d-01", year, int(month))
	dateFrom, err := time.Parse("2006-01-02", dateFromString)

	if err != nil {
		return nil, err
	}

	dateTo := dateFrom.AddDate(0, 1, 0)
	dateToString := dateTo.Format("2006-01-02")

	if err != nil {
		return nil, err
	}

	entities, err := repo.transactionDao.GetBreakDownByCategory(ctx, dateFromString, dateToString, user.Id)
	var totalAmount int64

	for _, e := range entities {
		totalAmount += e.Amount
	}

	for _, e := range entities {
		breakdown := domain.Breakdown{
			CategoryName: e.CategoryName,
			Amount:       money.New(e.Amount, user.Currency.Code),
			Percent:      float64(e.Amount) / float64(totalAmount) * 100,
		}
		breakdowns = append(breakdowns, breakdown)
	}

	return breakdowns, nil
}
