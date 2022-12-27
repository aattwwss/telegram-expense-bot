package dao

import (
	"context"

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
			SELECT id, datetime, category_id, description, user_id, amount, currency 
			FROM transaction 
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
		INSERT INTO transaction ( datetime, category_id, description, user_id, amount, currency)
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
			DELETE FROM expenditure_bot.transaction 
			WHERE id = $1 AND user_id = $2 ;
		`
	_, err := dao.db.Exec(ctx, sql, id, userId)
	if err != nil {
		return err
	}
	return nil
}

func (dao TransactionDAO) GetBreakDownByCategory(ctx context.Context, dateFrom string, dateTo string, userId int64) ([]entity.TransactionBreakdown, error) {
	entities := []entity.TransactionBreakdown{}
	sql := `
			SELECT c.name, SUM(amount) amount
			FROM expenditure_bot.transaction t
					 JOIN category c on t.category_id = c.id
			WHERE t.datetime >= $1
			  AND t.datetime < $2
			GROUP BY c.name
			ORDER BY amount DESC;	
		`
	err := pgxscan.Select(ctx, dao.db, &entities, sql, dateFrom, dateTo, userId)
	if err != nil {
		return nil, err
	}
	return entities, nil
}
