package repo

import (
	"context"
	"fmt"
	"math"
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
		Id:           t.Id,
		Datetime:     t.Datetime,
		CategoryId:   t.CategoryId,
		CategoryName: t.CategoryName,
		Description:  t.Description,
		UserId:       t.UserId,
		Amount:       t.Amount.Amount(),
		Currency:     t.Amount.Currency().Code,
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

func (repo TransactionRepo) GetTransactionBreakdownByCategory(ctx context.Context, month time.Month, year int, user domain.User) (domain.Breakdowns, *money.Money, error) {
	breakdowns := domain.Breakdowns{}

	dateFromString := fmt.Sprintf("%v-%02d-01", year, int(month))
	dateFrom, err := time.ParseInLocation("2006-01-02", dateFromString, user.Location)

	if err != nil {
		return nil, nil, err
	}

	dateTo := dateFrom.AddDate(0, 1, 0)

	if err != nil {
		return nil, nil, err
	}

	entities, err := repo.transactionDao.GetBreakdownByCategory(ctx, dateFrom, dateTo, user.Id)
	var totalAmount int64

	for _, e := range entities {
		totalAmount += e.Amount
	}

	for _, e := range entities {
		percent := float64(e.Amount) / float64(totalAmount) * 100
		breakdown := domain.Breakdown{
			CategoryName: e.CategoryName,
			Amount:       money.New(e.Amount, user.Currency.Code),
			Percent:      math.Round(percent*10) / 10,
		}
		breakdowns = append(breakdowns, breakdown)
	}

	return breakdowns, money.New(totalAmount, user.Currency.Code), nil
}

func (repo TransactionRepo) ListByMonthAndYear(ctx context.Context, month time.Month, year int, offset int, limit int, user domain.User) (domain.Transactions, int, error) {
	var transactions domain.Transactions

	dateFromString := fmt.Sprintf("%v-%02d-01", year, int(month))
	dateFrom, err := time.ParseInLocation("2006-01-02", dateFromString, user.Location)

	if err != nil {
		return nil, 0, err
	}

	dateTo := dateFrom.AddDate(0, 1, 0)

	totalCount, err := repo.transactionDao.CountListByMonthAndYear(ctx, dateFrom, dateTo, user.Id)
	if err != nil {
		return transactions, 0, err
	}
	if totalCount == 0 {
		return transactions, totalCount, nil
	}

	entities, err := repo.transactionDao.ListByMonthAndYear(ctx, dateFrom, dateTo, offset, limit, user.Id)
	if err != nil {
		return transactions, 0, err
	}

	for _, e := range entities {
		t := domain.Transaction{
			Id:           e.Id,
			Datetime:     e.Datetime,
			CategoryId:   e.CategoryId,
			CategoryName: e.CategoryName,
			Description:  e.Description,
			UserId:       e.UserId,
			Amount:       money.New(e.Amount, e.Currency),
		}
		transactions = append(transactions, t)
	}

	return transactions, totalCount, nil
}
