package dao

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/entity"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageContextDAO struct {
	db *pgxpool.Pool
}

func NewMessageContextDao(db *pgxpool.Pool) MessageContextDAO {
	return MessageContextDAO{db: db}
}

func (dao MessageContextDAO) Insert(ctx context.Context, messageContext entity.MessageContext) (int, error) {
	var lastInsertId int
	sql := `
		INSERT INTO message_context ( message )
		VALUES ($1) RETURNING id
		`
	err := dao.db.QueryRow(ctx, sql, messageContext.Message).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (dao MessageContextDAO) GetById(ctx context.Context, id int) (*entity.MessageContext, error) {
	var messageContextEntities []entity.MessageContext
	sql := `
			SELECT id, message
			FROM message_context
            WHERE id = $1;
			`
	err := pgxscan.Select(ctx, dao.db, &messageContextEntities, sql, id)
	if err != nil {
		return nil, err
	}
	return &messageContextEntities[0], nil
}

func (dao MessageContextDAO) DeleteById(ctx context.Context, id int) error {
	sql := `DELETE FROM message_context WHERE id = $1`
	_, err := dao.db.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	return nil
}
