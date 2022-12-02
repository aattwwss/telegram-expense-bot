package repo

import (
	"context"
	"github.com/aattwwss/telegram-expense-bot/dao"
	"github.com/aattwwss/telegram-expense-bot/domain"
)

type CategoryRepo struct {
	categoryDao dao.CategoryDAO
}

func NewCategoryRepo(categoryDao dao.CategoryDAO) CategoryRepo {
	return CategoryRepo{categoryDao: categoryDao}
}

func (repo CategoryRepo) FindByTransactionTypeId(ctx context.Context, transactionTypeId int64) ([]domain.Category, error) {
	var categories []domain.Category
	entities, err := repo.categoryDao.FindByTransactionTypeId(ctx, transactionTypeId)
	if err != nil {
		return nil, err
	}
	for _, c := range entities {
		category := domain.Category{
			Id:                c.Id,
			Name:              c.Name,
			TransactionTypeId: c.TransactionTypeId,
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (repo CategoryRepo) GetById(ctx context.Context, id int) (domain.Category, error) {
	c, err := repo.categoryDao.GetById(ctx, id)
	if err != nil {
		return domain.Category{}, err
	}
	category := domain.Category{
		Id:                c.Id,
		Name:              c.Name,
		TransactionTypeId: c.TransactionTypeId,
	}
	return category, nil
}
