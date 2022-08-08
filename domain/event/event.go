package event

import (
	"errors"
	"time"
)

type Event struct {
	// sport associated with this Event
	sport *Identifier
	// competition associated with this Event
	competition *Identifier
	// category associated with this Event
	category *Identifier
	// key-name tuple for the home team competitor of this Event
	home TeamIdentifier
	// key-name tuple for the away team competitor of this Event
	away TeamIdentifier
	// the event trading still active or not
	active bool
	// mapping between market key and all associated markets for this Event
	markets map[string]Market
	// name of this Event
	name string
	// slug for this Event
	key string
	// event cutoff time in string format "2006-01-02T15:04:05Z07:00" (RFC3339)
	cutoffTime time.Time
	// event start TRADING_LIVE time
	startTradingLiveTime time.Time
	// event inactive time
	inactiveTime time.Time
}

var (
	ErrMissingKey         = errors.New("missing event key")
	ErrMissingSport       = errors.New("missing sport info")
	ErrMissingCompetition = errors.New("missing competition info")
	ErrMissingCategory    = errors.New("missing category info")
	ErrMissingCutOffTime  = errors.New("missing cut off time")
)

func NewEvent(sport *Identifier, competition *Identifier, category *Identifier, home TeamIdentifier, away TeamIdentifier, active bool, markets map[string]Market, name string, key string, cutOffTime time.Time) (*Event, error) {
	if key == "" {
		return nil, ErrMissingKey
	}
	if sport == nil {
		return nil, ErrMissingSport
	}
	if competition == nil {
		return nil, ErrMissingCompetition
	}
	if category == nil {
		return nil, ErrMissingCategory
	}
	if cutOffTime.IsZero() {
		return nil, ErrMissingCutOffTime
	}

	var startTradingLiveTime time.Time
	if active {
		startTradingLiveTime = time.Now()
	}

	return &Event{
		sport:                sport,
		competition:          competition,
		category:             category,
		home:                 home,
		away:                 away,
		active:               active,
		markets:              markets,
		name:                 name,
		key:                  key,
		cutoffTime:           cutOffTime,
		startTradingLiveTime: startTradingLiveTime,
	}, nil
}

// read only properties
func (e Event) Sport() Identifier {
	return *e.sport
}

func (e Event) Category() Identifier {
	return *e.category
}

func (e Event) Competition() Identifier {
	return *e.competition
}

func (e Event) Home() TeamIdentifier {
	return e.home
}

func (e Event) Away() TeamIdentifier {
	return e.away
}

func (e Event) Active() bool {
	return e.active
}

func (e Event) Market() map[string]Market {
	return e.markets
}

func (e *Event) Inactivate() {
	e.active = false
	e.inactiveTime = time.Now()
}

func (e Event) Key() string {
	return e.key
}

func (e Event) Name() string {
	return e.name
}

func (e Event) CutOffTime() time.Time {
	return e.cutoffTime
}

func (e Event) StartTradingLiveTime() time.Time {
	return e.startTradingLiveTime
}

func (e *Event) SetstartTradingLiveTime(time time.Time) {
	e.startTradingLiveTime = time
}

func (e Event) InactiveTime() time.Time {
	return e.inactiveTime
}
