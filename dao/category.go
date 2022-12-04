package dao

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/entity"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryDAO struct {
	db *pgxpool.Pool
}

func NewCategoryDAO(db *pgxpool.Pool) CategoryDAO {
	return CategoryDAO{db: db}
}

func (dao CategoryDAO) FindByTransactionTypeId(ctx context.Context, transactionTypeId int64) ([]*entity.Category, error) {
	var categories []*entity.Category
	sql := `
			SELECT id, name, transaction_type_id 
			FROM category 
			WHERE transaction_type_id = $1
			ORDER BY display_order
			`
	err := pgxscan.Select(ctx, dao.db, &categories, sql, transactionTypeId)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (dao CategoryDAO) GetById(ctx context.Context, id int) (entity.Category, error) {
	var categories []*entity.Category
	sql := `
			SELECT id, name, transaction_type_id 
			FROM category 
			WHERE id = $1
			`
	err := pgxscan.Select(ctx, dao.db, &categories, sql, id)
	if err != nil {
		return entity.Category{}, err
	}
	return *categories[0], nil
}
