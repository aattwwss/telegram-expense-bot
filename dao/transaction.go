package dao

import (
	"context"
	"time"

	"github.com/aattwwss/telegram-expense-bot/entity"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionDAO struct {
	db *pgxpool.Pool
}

func NewTransactionDao(db *pgxpool.Pool) TransactionDAO {
	return TransactionDAO{db: db}
}

func (dao TransactionDAO) GetById(ctx context.Context, id int, userId int64) (entity.Transaction, error) {
	var transactions []*entity.Transaction
	sql := `
			SELECT t.id, t.datetime, t.category_id, t.description, t.user_id, t.amount, t.currency, c.name as category_name
			FROM transaction t JOIN category c on t.category_id = c.id
			WHERE id = $1 and user_id = $2
			`
	err := pgxscan.Select(ctx, dao.db, &transactions, sql, id, userId)
	if err != nil {
		return entity.Transaction{}, err
	}
	return *transactions[0], nil
}

func (dao TransactionDAO) FindLatestByUserId(ctx context.Context, userId int64) (*entity.Transaction, error) {
	var transactions []*entity.Transaction
	sql := `
			SELECT id, datetime, category_id, description, user_id, amount, currency 
			FROM transaction 
			WHERE user_id = $1
			ORDER BY datetime DESC LIMIT 1;
			`
	err := pgxscan.Select(ctx, dao.db, &transactions, sql, userId)
	if err != nil {
		return nil, err
	}
	if len(transactions) == 0 {
		return nil, nil
	}
	return transactions[0], nil

}

func (dao TransactionDAO) Insert(ctx context.Context, transaction entity.Transaction) error {
	sql := `
		INSERT INTO transaction (datetime, category_id, description, user_id, amount, currency)
		VALUES ($1, $2, $3, $4, $5, $6)
		`
	_, err := dao.db.Exec(ctx, sql, transaction.Datetime, transaction.CategoryId, transaction.Description, transaction.UserId, transaction.Amount, transaction.Currency)
	if err != nil {
		return err
	}
	return nil
}

func (dao TransactionDAO) DeleteById(ctx context.Context, id int, userId int64) error {
	sql := `
			DELETE FROM transaction 
			WHERE id = $1 AND user_id = $2 ;
		`
	_, err := dao.db.Exec(ctx, sql, id, userId)
	if err != nil {
		return err
	}
	return nil
}

func (dao TransactionDAO) GetBreakdownByCategory(ctx context.Context, dateFrom time.Time, dateTo time.Time, userId int64) ([]entity.TransactionBreakdown, error) {
	var entities []entity.TransactionBreakdown
	sql := `
			SELECT c.name as     category_name,
			       sum(t.amount) amount
			FROM transaction t JOIN category c on t.category_id = c.id
			WHERE datetime >= $1::timestamptz
			AND datetime < $2::timestamptz
			AND user_id = $3
			GROUP BY c.name
			ORDER BY amount DESC;
		`
	err := pgxscan.Select(ctx, dao.db, &entities, sql, dateFrom.Format(time.RFC3339), dateTo.Format(time.RFC3339), userId)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (dao TransactionDAO) ListByMonthAndYear(ctx context.Context, dateFrom time.Time, dateTo time.Time, offset int, limit int, userId int64) ([]entity.Transaction, error) {
	var entities []entity.Transaction
	sql := `
			SELECT t.id, t.datetime, t.category_id, t.description, t.user_id, t.amount, t.currency, c.name as category_name
			FROM transaction t JOIN category c on t.category_id = c.id
		    WHERE t.datetime >= $1::timestamptz
			  AND t.datetime < $2::timestamptz
			  AND t.user_id = $3
		    ORDER BY t.datetime DESC
			OFFSET $4 LIMIT $5
		`
	err := pgxscan.Select(ctx, dao.db, &entities, sql, dateFrom.Format(time.RFC3339), dateTo.Format(time.RFC3339), userId, offset, limit)
	if err != nil {
		return nil, err
	}
	return entities, nil
}
