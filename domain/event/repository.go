package event

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, event Event) error
	ListEvents(ctx context.Context, first int, page int, sportKey string, categoryKey string, competitionKey string) ([]Event, error)
	GetEvent(ctx context.Context, eventKey string) (Event, error)
	ListEventsCutOffSoon(ctx context.Context) ([]Event, error)
}
