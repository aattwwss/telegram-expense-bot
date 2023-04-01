package repo

import (
	"context"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/domain"
	"github.com/aattwwss/telegram-expense-bot/entity"
)

type UserRepo struct {
	userDao dao.UserDAO
}

func NewUserRepo(userDao dao.UserDAO) UserRepo {
	return UserRepo{userDao: userDao}
}

func (repo UserRepo) FindUserById(ctx context.Context, id int64) (*domain.User, error) {
	userEntity, err := repo.userDao.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	if userEntity == nil {
		return nil, nil
	}

	loc, err := time.LoadLocation(userEntity.Timezone)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Id:       userEntity.Id,
		Locale:   userEntity.Locale,
		Currency: money.GetCurrency(userEntity.Currency),
		Location: loc,
	}

	return &user, nil
}

func (repo UserRepo) Add(ctx context.Context, user domain.User) error {
	userEntity := entity.User{
		Id:       user.Id,
		Locale:   user.Locale,
		Currency: user.Currency.Code,
		Timezone: user.Location.String(),
	}
	err := repo.userDao.Insert(ctx, userEntity)
	if err != nil {
		return err
	}
	return nil
}

func (repo UserRepo) Update(ctx context.Context, user domain.User) error {
	userEntity := entity.User{
		Id:       user.Id,
		Locale:   user.Locale,
		Currency: user.Currency.Code,
		Timezone: user.Location.String(),
	}
	err := repo.userDao.Update(ctx, userEntity)
	if err != nil {
		return err
	}
	return nil
}
