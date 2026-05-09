package handler

import (
	"context"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/entity"
)

type UserRepo interface {
	FindUserById(ctx context.Context, id int64) (*domain.User, error)
	Add(ctx context.Context, user domain.User) error
}

type TransactionRepo interface {
	Add(ctx context.Context, t domain.Transaction) error
	GetById(ctx context.Context, id int, userId int64) (domain.Transaction, error)
	FindLastestByUserId(ctx context.Context, userId int64) (*domain.Transaction, error)
	DeleteById(ctx context.Context, id int, userId int64) error
	GetTransactionBreakdownByCategory(ctx context.Context, month time.Month, year int, user domain.User) (domain.Breakdowns, *money.Money, error)
	ListByMonthAndYear(ctx context.Context, q entity.TransactionListQuery) (domain.Transactions, int, error)
}

type MessageContextRepo interface {
	Add(ctx context.Context, chatId int64, messageId int, message string) (int, error)
	GetMessageById(ctx context.Context, id int) (string, error)
	DeleteById(ctx context.Context, id int) error
}

type TransactionTypeRepo interface {
	GetById(ctx context.Context, id int) (*entity.TransactionType, error)
}

type CategoryRepo interface {
	FindAll(ctx context.Context) ([]*entity.Category, error)
	GetById(ctx context.Context, id int) (*entity.Category, error)
}
