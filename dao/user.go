package dao

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/entity"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDAO struct {
	db *pgxpool.Pool
}

func NewUserDao(db *pgxpool.Pool) UserDAO {
	return UserDAO{db: db}
}

func (dao UserDAO) FindUserById(ctx context.Context, id int64) (*entity.User, error) {
	var users []*entity.User
	err := pgxscan.Select(ctx, dao.db, &users, `SELECT id, first_name, last_name, username FROM telegram_user WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	if len(users) != 1 {
		return nil, nil
	}
	return users[0], nil
}

func (dao UserDAO) Insert(ctx context.Context, user entity.User) error {
	sql := `
		INSERT INTO telegram_user (id, is_bot,first_name, last_name, username)
		VALUES ($1, $2, $3, $4, $5)
		`
	_, err := dao.db.Exec(ctx, sql, user.Id, user.IsBot, user.FirstName, user.LastName, user.Username)
	if err != nil {
		return err
	}
	return nil
}
