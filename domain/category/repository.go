package category

import (
	"context"
)

type Repository interface {
	ListCategories(ctx context.Context, first int, page int, sportKey string) ([]Category, error)
	GetCategory(ctx context.Context, categoryKey string) (Category, error)
}
