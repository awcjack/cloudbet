package query

import (
	"context"

	"github.com/awcjack/cloudbet/domain/event"
)

type ListEventsHandler struct {
	eventRepo event.Repository
	logger    logger
}

func NewListEventsHandler(eventRepo event.Repository, logger logger) *ListEventsHandler {
	return &ListEventsHandler{
		eventRepo: eventRepo,
		logger:    logger,
	}
}

func (l ListEventsHandler) Handle(ctx context.Context, first int, page int, sportKey string, categoryKey string, competitionKey string) ([]event.Event, error) {
	return l.eventRepo.ListEvents(ctx, first, page, sportKey, categoryKey, competitionKey)
}

type GetEventHandler struct {
	eventRepo event.Repository
	logger    logger
}

func NewGetEventHandler(eventRepo event.Repository, logger logger) *GetEventHandler {
	return &GetEventHandler{
		eventRepo: eventRepo,
		logger:    logger,
	}
}

func (g GetEventHandler) Handle(ctx context.Context, eventKey string) (event.Event, error) {
	return g.eventRepo.GetEvent(ctx, eventKey)
}
