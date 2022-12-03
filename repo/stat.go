package repo

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"time"
)

type StatRepo struct {
	statDao dao.StatDAO
}

type GetMonthlySearchParam struct {
	Location           time.Location
	PrevMonthIntervals int
	UserId             int64
}

func NewStatRepo(statDao dao.StatDAO) StatRepo {
	return StatRepo{statDao: statDao}
}

func (repo StatRepo) GetMonthly(ctx context.Context, param GetMonthlySearchParam) (domain.MonthlySummaries, error) {
	var summaries domain.MonthlySummaries
	entities, err := repo.statDao.GetMonthly(ctx, param.PrevMonthIntervals, param.UserId)
	if err != nil {
		return nil, err
	}

	for _, entity := range entities {
		summary := domain.MonthlySummary{
			Month:                entity.Datetime.Month(),
			Year:                 entity.Datetime.Year(),
			Amount:               entity.Amount,
			TransactionTypeLabel: entity.TransactionTypeLabel,
			Multiplier:           entity.Multiplier,
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}
