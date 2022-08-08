package sport

import (
	"context"
)

type Repository interface {
	ListSports(ctx context.Context, first int, page int) ([]Sport, error)
	GetSport(ctx context.Context, sportKey string) (Sport, error)
}
