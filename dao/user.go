package dao

import (
	"context"
	"errors"

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
	err := pgxscan.Select(ctx, dao.db, &users, `SELECT id, locale, currency, timezone, current_context FROM app_user WHERE id = $1`, id)
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
		INSERT INTO app_user (id, locale, currency, timezone, current_context)
		VALUES ($1, $2, $3, $4, $5)
		`
	_, err := dao.db.Exec(ctx, sql, user.Id, user.Locale, user.Currency, user.Timezone, user.CurrentContext)
	if err != nil {
		return err
	}
	return nil
}

func (dao UserDAO) Update(ctx context.Context, user entity.User) error {
	if user.Id == 0 {
		return errors.New("user id cannot be 0 or empty")
	}

	sql := `
		UPDATE app_user SET 
	      locale = $2, 
	      currency = $3, 
	      timezone = $4,
	      current_context = $5
	    WHERE id = $1
		`
	_, err := dao.db.Exec(ctx, sql, user.Id, user.Locale, user.Currency, user.Timezone, user.CurrentContext)
	if err != nil {
		return err
	}
	return nil
}
