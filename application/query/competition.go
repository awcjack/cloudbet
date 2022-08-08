package query

import (
	"context"

	"github.com/awcjack/cloudbet/domain/competition"
)

type ListCompetitionsHandler struct {
	competitionRepo competition.Repository
	logger          logger
}

func NewListCompetitionHandler(competitionRepo competition.Repository, logger logger) *ListCompetitionsHandler {
	return &ListCompetitionsHandler{
		competitionRepo: competitionRepo,
		logger:          logger,
	}
}

func (l ListCompetitionsHandler) Handle(ctx context.Context, first int, page int, sportKey string) ([]competition.Competition, error) {
	return l.competitionRepo.ListCompetitions(ctx, first, page, sportKey)
}

type GetCompetitionHandler struct {
	competitionRepo competition.Repository
	logger          logger
}

func NewGetCompetitionHandler(competitionRepo competition.Repository, logger logger) *GetCompetitionHandler {
	return &GetCompetitionHandler{
		competitionRepo: competitionRepo,
		logger:          logger,
	}
}

func (g GetCompetitionHandler) Handle(ctx context.Context, competitionKey string) (competition.Competition, error) {
	return g.competitionRepo.GetCompetition(ctx, competitionKey)
}
