package query

import (
	"context"

	"github.com/awcjack/cloudbet/domain/category"
)

type ListCategoriesHandler struct {
	categoryRepo category.Repository
	logger       logger
}

func NewListCategoriesHandler(categoryRepo category.Repository, logger logger) *ListCategoriesHandler {
	return &ListCategoriesHandler{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (l ListCategoriesHandler) Handle(ctx context.Context, first int, page int, sportKey string) ([]category.Category, error) {
	return l.categoryRepo.ListCategories(ctx, first, page, sportKey)
}

type GetCategoryHandler struct {
	categoryRepo category.Repository
	logger       logger
}

func NewGetCategoryHandler(categoryRepo category.Repository, logger logger) *GetCategoryHandler {
	return &GetCategoryHandler{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (g GetCategoryHandler) Handle(ctx context.Context, categoryKey string) (category.Category, error) {
	return g.categoryRepo.GetCategory(ctx, categoryKey)
}
