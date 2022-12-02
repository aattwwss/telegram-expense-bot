package dao

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/entity"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionTypeDAO struct {
	db *pgxpool.Pool
}

func NewTransactionTypeDAO(db *pgxpool.Pool) TransactionTypeDAO {
	return TransactionTypeDAO{db: db}
}

func (dao TransactionTypeDAO) GetAll(ctx context.Context) ([]*entity.TransactionType, error) {
	var types []*entity.TransactionType
	sql := `
			SELECT id, name, multiplier, reply_text
			FROM transaction_type
			ORDER BY display_order
			`
	err := pgxscan.Select(ctx, dao.db, &types, sql)
	if err != nil {
		return nil, err
	}
	return types, nil
}

func (dao TransactionTypeDAO) GetById(ctx context.Context, id int64) (*entity.TransactionType, error) {
	var types []*entity.TransactionType
	sql := `
			SELECT id, name, multiplier, reply_text
			FROM transaction_type
            WHERE id = $1
			ORDER BY display_order
			`
	err := pgxscan.Select(ctx, dao.db, &types, sql, id)
	if err != nil {
		return nil, err
	}
	return types[0], nil
}
