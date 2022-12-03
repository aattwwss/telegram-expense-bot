package dao

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/entity"
	"github.com/aattwwss/telegram-expense-bot/util"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatDAO struct {
	db *pgxpool.Pool
}

func NewStatDAO(db *pgxpool.Pool) StatDAO {
	return StatDAO{db: db}
}

func (dao StatDAO) GetMonthly(ctx context.Context, from util.YearMonth, to util.YearMonth, userId int64) ([]*entity.MonthlySummary, error) {
	fromString, err := from.String("2006-01")
	if err != nil {
		return nil, err
	}

	toString, err := to.String("2006-01")
	if err != nil {
		return nil, err
	}

	var summaries []*entity.MonthlySummary

	sql := `
			SELECT datetime, amount, transaction_type_label, multiplier
			FROM monthly_transaction_agg
			WHERE datetime >= TO_TIMESTAMP($1, 'YYYY-MM') 
			  AND datetime <= TO_TIMESTAMP($2, 'YYYY-MM') 
			  AND user_id = $3
			`

	err = pgxscan.Select(ctx, dao.db, &summaries, sql, fromString, toString, userId)
	if err != nil {
		return nil, err
	}
	return summaries, nil
}
