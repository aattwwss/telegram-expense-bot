package dao

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionDAO struct {
	db *pgxpool.Pool
}

func NewTransactionDao(db *pgxpool.Pool) TransactionDAO {
	return TransactionDAO{db: db}
}

func (dao TransactionDAO) Insert(ctx context.Context, transaction entity.Transaction) error {
	sql := `
		INSERT INTO transaction ( datetime, category_id, description, user_id, amount, currency)
		VALUES ($1, $2, $3, $4, $5, $6)
		`
	_, err := dao.db.Exec(ctx, sql, transaction.Datetime, transaction.CategoryId, transaction.Description, transaction.UserId, transaction.Amount, transaction.Currency)
	if err != nil {
		return err
	}
	return nil
}
