package query

import (
	"context"

	"github.com/awcjack/cloudbet/domain/sport"
)

type ListSportsHandler struct {
	sportRepo sport.Repository
	logger    logger
}

func NewListSportHandler(sportRepo sport.Repository, logger logger) *ListSportsHandler {
	return &ListSportsHandler{
		sportRepo: sportRepo,
		logger:    logger,
	}
}

func (l ListSportsHandler) Handle(ctx context.Context, first int, page int) ([]sport.Sport, error) {
	return l.sportRepo.ListSports(ctx, first, page)
}

type GetSportHandler struct {
	sportRepo sport.Repository
	logger    logger
}

func NewGetSportHandler(sportRepo sport.Repository, logger logger) *GetSportHandler {
	return &GetSportHandler{
		sportRepo: sportRepo,
		logger:    logger,
	}
}

func (g GetSportHandler) Handle(ctx context.Context, sportKey string) (sport.Sport, error) {
	return g.sportRepo.GetSport(ctx, sportKey)
}
