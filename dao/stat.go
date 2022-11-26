package dao

import (
	"context"
	"fmt"
	"github.com/aattwwss/telegram-expense-bot/entity"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type StatDAO struct {
	db *pgxpool.Pool
}

type GetMonthlySearchParam struct {
	Location           time.Location
	PrevMonthIntervals int
	UserId             int64
}

func NewStatDAO(db *pgxpool.Pool) StatDAO {
	return StatDAO{db: db}
}

func (dao StatDAO) GetMonthly(ctx context.Context, param GetMonthlySearchParam) ([]*entity.MonthlySummary, error) {
	var summaries []*entity.MonthlySummary
	sql := `
			select date_part('month', datetime) as month,
				   date_part('year', datetime)  as year,
				   amount                       as amount,
				   transaction_type_label       as transaction_type_label,
				   multiplier                   as multiplier
			from (SELECT date_trunc('month', datetime AT time zone '%[1]s') as datetime,
						 sum(amount)                                        as amount,
						 tt.name                                            as transaction_type_label,
						 tt.multiplier                                      as multiplier
				  FROM expenditure_bot.transaction t
						   join expenditure_bot.category c on c.id = t.category_id
						   join expenditure_bot.transaction_type tt on c.transaction_type_id = tt.id
				  WHERE datetime > CURRENT_TIMESTAMP - $1 * interval '1 month'
					AND t.user_id = $2
				  GROUP BY tt.id, date_trunc('month', datetime AT time zone '%[1]s')
				  ORDER BY date_trunc('month', datetime AT time zone '%[1]s') ASC) as a;
		`

	sql = fmt.Sprintf(sql, param.Location.String())
	err := pgxscan.Select(ctx, dao.db, &summaries, sql, param.PrevMonthIntervals, param.UserId)
	if err != nil {
		return nil, err
	}
	return summaries, nil
}
