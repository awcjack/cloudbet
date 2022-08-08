package competition

import (
	"context"
)

type Repository interface {
	ListCompetitions(ctx context.Context, first int, page int, sportKey string) ([]Competition, error)
	GetCompetition(ctx context.Context, competitionKey string) (Competition, error)
}
