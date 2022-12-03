package dao

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/entity"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatDAO struct {
	db *pgxpool.Pool
}

func NewStatDAO(db *pgxpool.Pool) StatDAO {
	return StatDAO{db: db}
}

func (dao StatDAO) GetMonthly(ctx context.Context, prevMonths int, userId int64) ([]*entity.MonthlySummary, error) {
	var summaries []*entity.MonthlySummary

	sql := `
			SELECT datetime, amount, transaction_type_label, multiplier
			FROM monthly_transaction_agg
			WHERE datetime > CURRENT_TIMESTAMP - $1 * interval '1 month'
			  AND user_id = $2
			`

	err := pgxscan.Select(ctx, dao.db, &summaries, sql, prevMonths, userId)
	if err != nil {
		return nil, err
	}
	return summaries, nil
}
