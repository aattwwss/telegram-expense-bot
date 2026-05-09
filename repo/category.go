package repo

import (
	"context"

	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/entity"
)

type CategoryRepo struct {
	categoryDao dao.CategoryDAO
}

func NewCategoryRepo(categoryDao dao.CategoryDAO) CategoryRepo {
	return CategoryRepo{categoryDao: categoryDao}
}

func (repo CategoryRepo) FindAll(ctx context.Context) ([]*entity.Category, error) {
	return repo.categoryDao.FindAll(ctx)
}

func (repo CategoryRepo) FindByTransactionTypeId(ctx context.Context, transactionTypeId int) ([]*entity.Category, error) {
	return repo.categoryDao.FindByTransactionTypeId(ctx, transactionTypeId)
}

func (repo CategoryRepo) GetById(ctx context.Context, id int) (*entity.Category, error) {
	e, err := repo.categoryDao.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
