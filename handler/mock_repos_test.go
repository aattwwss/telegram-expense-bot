package handler

import (
	"context"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/entity"
)

type mockUserRepo struct {
	findByIdFn func(ctx context.Context, id int64) (*domain.User, error)
	addFn      func(ctx context.Context, user domain.User) error
}

func (m mockUserRepo) FindUserById(ctx context.Context, id int64) (*domain.User, error) {
	return m.findByIdFn(ctx, id)
}

func (m mockUserRepo) Add(ctx context.Context, user domain.User) error {
	return m.addFn(ctx, user)
}

type mockTransactionRepo struct {
	addFn                          func(ctx context.Context, t domain.Transaction) error
	getByIdFn                      func(ctx context.Context, id int, userId int64) (domain.Transaction, error)
	findLatestByUserIdFn           func(ctx context.Context, userId int64) (*domain.Transaction, error)
	deleteByIdFn                   func(ctx context.Context, id int, userId int64) error
	getTransactionBreakdownByCatFn func(ctx context.Context, month time.Month, year int, user domain.User) (domain.Breakdowns, *money.Money, error)
	listByMonthAndYearFn           func(ctx context.Context, q entity.TransactionListQuery) (domain.Transactions, int, error)
}

func (m mockTransactionRepo) Add(ctx context.Context, t domain.Transaction) error {
	return m.addFn(ctx, t)
}

func (m mockTransactionRepo) GetById(ctx context.Context, id int, userId int64) (domain.Transaction, error) {
	return m.getByIdFn(ctx, id, userId)
}

func (m mockTransactionRepo) FindLastestByUserId(ctx context.Context, userId int64) (*domain.Transaction, error) {
	return m.findLatestByUserIdFn(ctx, userId)
}

func (m mockTransactionRepo) DeleteById(ctx context.Context, id int, userId int64) error {
	return m.deleteByIdFn(ctx, id, userId)
}

func (m mockTransactionRepo) GetTransactionBreakdownByCategory(ctx context.Context, month time.Month, year int, user domain.User) (domain.Breakdowns, *money.Money, error) {
	return m.getTransactionBreakdownByCatFn(ctx, month, year, user)
}

func (m mockTransactionRepo) ListByMonthAndYear(ctx context.Context, q entity.TransactionListQuery) (domain.Transactions, int, error) {
	return m.listByMonthAndYearFn(ctx, q)
}

type mockMessageContextRepo struct {
	addFn        func(ctx context.Context, chatId int64, messageId int, message string) (int, error)
	getMsgByIdFn func(ctx context.Context, id int) (string, error)
	deleteByIdFn func(ctx context.Context, id int) error
}

func (m mockMessageContextRepo) Add(ctx context.Context, chatId int64, messageId int, message string) (int, error) {
	return m.addFn(ctx, chatId, messageId, message)
}

func (m mockMessageContextRepo) GetMessageById(ctx context.Context, id int) (string, error) {
	return m.getMsgByIdFn(ctx, id)
}

func (m mockMessageContextRepo) DeleteById(ctx context.Context, id int) error {
	return m.deleteByIdFn(ctx, id)
}

type mockTransactionTypeRepo struct {
	getByIdFn func(ctx context.Context, id int) (*entity.TransactionType, error)
}

func (m mockTransactionTypeRepo) GetById(ctx context.Context, id int) (*entity.TransactionType, error) {
	return m.getByIdFn(ctx, id)
}

type mockCategoryRepo struct {
	findAllFn func(ctx context.Context) ([]*entity.Category, error)
	getByIdFn func(ctx context.Context, id int) (*entity.Category, error)
}

func (m mockCategoryRepo) FindAll(ctx context.Context) ([]*entity.Category, error) {
	return m.findAllFn(ctx)
}

func (m mockCategoryRepo) GetById(ctx context.Context, id int) (*entity.Category, error) {
	return m.getByIdFn(ctx, id)
}
