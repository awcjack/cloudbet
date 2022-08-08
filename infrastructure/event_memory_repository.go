package infrastructure

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/awcjack/cloudbet/domain/category"
	"github.com/awcjack/cloudbet/domain/competition"
	"github.com/awcjack/cloudbet/domain/event"
	"github.com/awcjack/cloudbet/domain/sport"
)

var (
	ErrOutOfRange    = errors.New("page exceed max page")
	ErrEventNotFound = errors.New("event not found")
)

type MemoryRepository struct {
	// TODO: change to map?
	sports             []sport.Sport
	categories         []category.Category
	sportsCategories   map[string][]string
	competitions       []competition.Competition
	sportsCompetitions map[string][]string
	events             []event.Event
	sportsEvents       map[string][]string
	categoriesEvents   map[string][]string
	competitionsEvents map[string][]string
	lock               *sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		sports:             make([]sport.Sport, 0),
		categories:         make([]category.Category, 0),
		sportsCategories:   make(map[string][]string),
		competitions:       make([]competition.Competition, 0),
		sportsCompetitions: make(map[string][]string),
		events:             make([]event.Event, 0),
		sportsEvents:       make(map[string][]string),
		categoriesEvents:   make(map[string][]string),
		competitionsEvents: make(map[string][]string),
		lock:               &sync.RWMutex{},
	}
}

func (m *MemoryRepository) Save(_ context.Context, event event.Event) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	sportKey := event.Sport().Key()
	categoryKey := event.Category().Key()
	competitionKey := event.Competition().Key()
	eventKey := event.Key()
	m.sports = append(m.sports, sport.NewSport(event.Sport().Name(), sportKey, 0))
	found := false
	for _, v := range m.sports {
		if v.Key() == sportKey {
			found = true
			break
		}
	}
	if !found {
		m.sports = append(m.sports, sport.NewSport(event.Sport().Name(), sportKey, 0))
	}
	found = false
	for _, v := range m.categories {
		if v.Key() == categoryKey {
			found = true
			break
		}
	}
	if !found {
		m.categories = append(m.categories, category.NewCategory(event.Category().Name(), categoryKey))
	}
	_, ok := m.sportsCategories[sportKey]
	if ok {
		found := false
		for _, v := range m.sportsCategories[sportKey] {
			if v == categoryKey {
				found = true
				break
			}
		}
		if !found {
			m.sportsCategories[sportKey] = append(m.sportsCategories[sportKey], categoryKey)
		}
	} else {
		m.sportsCategories[sportKey] = append(m.sportsCategories[sportKey], categoryKey)
	}
	found = false
	for _, v := range m.competitions {
		if v.Key() == competitionKey {
			found = true
			break
		}
	}
	if !found {
		m.competitions = append(m.competitions, competition.NewCompetition(event.Competition().Name(), competitionKey))
	}
	_, ok = m.sportsCompetitions[sportKey]
	if ok {
		found := false
		for _, v := range m.sportsCompetitions[sportKey] {
			if v == competitionKey {
				found = true
				break
			}
		}
		if !found {
			m.sportsCompetitions[sportKey] = append(m.sportsCompetitions[sportKey], competitionKey)
		}
	} else {
		m.sportsCompetitions[sportKey] = append(m.sportsCompetitions[sportKey], competitionKey)
	}
	found = false
	for _, v := range m.events {
		if v.Key() == eventKey {
			found = true
			break
		}
	}
	if !found {
		m.events = append(m.events, event)
	}
	_, ok = m.categoriesEvents[categoryKey]
	if ok {
		found := false
		for _, v := range m.sportsEvents[sportKey] {
			if v == eventKey {
				found = true
				break
			}
		}
		if !found {
			m.sportsEvents[sportKey] = append(m.sportsEvents[sportKey], eventKey)
		}
	} else {
		m.sportsEvents[sportKey] = append(m.sportsEvents[sportKey], eventKey)
	}
	_, ok = m.categoriesEvents[categoryKey]
	if ok {
		found := false
		for _, v := range m.categoriesEvents[categoryKey] {
			if v == eventKey {
				found = true
				break
			}
		}
		if !found {
			m.categoriesEvents[categoryKey] = append(m.categoriesEvents[categoryKey], eventKey)
		}
	} else {
		m.categoriesEvents[categoryKey] = append(m.categoriesEvents[categoryKey], eventKey)
	}
	_, ok = m.competitionsEvents[competitionKey]
	if ok {
		found := false
		for _, v := range m.competitionsEvents[competitionKey] {
			if v == eventKey {
				found = true
				break
			}
		}
		if !found {
			m.competitionsEvents[competitionKey] = append(m.competitionsEvents[competitionKey], eventKey)
		}
	} else {
		m.competitionsEvents[competitionKey] = append(m.competitionsEvents[competitionKey], eventKey)
	}

	return nil
}

func (m MemoryRepository) ListSports(_ context.Context, first int, page int) ([]sport.Sport, error) {
	if len(m.sports) == 0 {
		return nil, ErrEventNotFound
	}
	if len(m.sports) < (page-1)*first {
		return nil, ErrOutOfRange
	}

	var target []sport.Sport
	if len(m.sports) <= (page)*first {
		target = m.sports[(page-1)*first : len(m.sports)]
	} else {
		target = m.sports[(page-1)*first : (page)*first]
	}

	for i, v := range target {
		counter := 0
		var timer int64 = 0
		for _, eventValue := range m.events {
			if v.Key() == eventValue.Sport().Key() && !eventValue.StartTradingLiveTime().IsZero() && !eventValue.InactiveTime().IsZero() {
				timer += int64(eventValue.InactiveTime().Sub(eventValue.StartTradingLiveTime()) / time.Millisecond)
				counter++
			}
		}
		target[i].SetLiveTime(float64(timer) / float64(counter))
	}
	return target, nil
}

func (m MemoryRepository) GetSport(_ context.Context, sportKey string) (sport.Sport, error) {
	var found bool
	var targetSport sport.Sport
	for _, v := range m.sports {
		if v.Key() == sportKey {
			found = true
			targetSport = v
			break
		}
	}
	if !found {
		return sport.Sport{}, ErrEventNotFound
	}

	counter := 0
	var timer int64 = 0
	for _, eventValue := range m.events {
		if sportKey == eventValue.Sport().Key() && !eventValue.StartTradingLiveTime().IsZero() && !eventValue.InactiveTime().IsZero() {
			timer += int64(eventValue.InactiveTime().Sub(eventValue.StartTradingLiveTime()) / time.Millisecond)
			counter++
		}
	}

	return sport.NewSport(targetSport.Name(), sportKey, float64(timer)/float64(counter)), nil
}

func (m MemoryRepository) ListCategories(_ context.Context, first int, page int, sportKey string) ([]category.Category, error) {
	if len(m.categories) == 0 {
		return nil, ErrEventNotFound
	}
	if sportKey == "" {
		if len(m.categories) < (page-1)*first {
			return nil, ErrOutOfRange
		}

		var target []category.Category
		if len(m.categories) <= (page)*first {
			target = m.categories[(page-1)*first:]
		} else {
			target = m.categories[(page-1)*first : (page)*first]
		}

		return target, nil
	}

	v, ok := m.sportsCategories[sportKey]
	if !ok {
		return nil, ErrEventNotFound
	}

	if len(v) < (page-1)*first {
		return nil, ErrOutOfRange
	}
	var targets []string
	if len(v) <= (page)*first {
		targets = v[(page-1)*first:]
	} else {
		targets = v[(page-1)*first : (page)*first]
	}
	result := make([]category.Category, 0, len(targets))
	for _, target := range targets {
		for _, v := range m.categories {
			if target == v.Key() {
				result = append(result, v)
			}
		}
	}
	return result, nil
}

func (m MemoryRepository) GetCategory(_ context.Context, categoryKey string) (category.Category, error) {
	var found bool
	var targetCategory category.Category
	for _, v := range m.categories {
		if v.Key() == categoryKey {
			found = true
			targetCategory = v
			break
		}
	}
	if !found {
		return category.Category{}, ErrEventNotFound
	}

	return targetCategory, nil
}

func (m MemoryRepository) ListCompetitions(_ context.Context, first int, page int, sportKey string) ([]competition.Competition, error) {
	if len(m.competitions) == 0 {
		return nil, ErrEventNotFound
	}
	if sportKey == "" {
		if len(m.competitions) < (page-1)*first {
			return nil, ErrOutOfRange
		}

		var target []competition.Competition
		if len(m.competitions) <= (page)*first {
			target = m.competitions[(page-1)*first : len(m.competitions)]
		} else {
			target = m.competitions[(page-1)*first : (page)*first]
		}

		return target, nil
	}

	v, ok := m.sportsCompetitions[sportKey]
	if !ok {
		return nil, ErrEventNotFound
	}

	if len(v) < (page-1)*first {
		return nil, ErrOutOfRange
	}

	var targets []string
	if len(v) <= (page)*first {
		targets = v[(page-1)*first:]
	} else {
		targets = v[(page-1)*first : (page)*first]
	}
	result := make([]competition.Competition, 0, len(targets))
	for _, target := range targets {
		for _, v := range m.competitions {
			if target == v.Key() {
				result = append(result, v)
			}
		}
	}
	return result, nil
}

func (m MemoryRepository) GetCompetition(_ context.Context, competitionKey string) (competition.Competition, error) {
	var found bool
	var targetCompetition competition.Competition
	for _, v := range m.competitions {
		if v.Key() == competitionKey {
			found = true
			targetCompetition = v
			break
		}
	}
	if !found {
		return competition.Competition{}, ErrEventNotFound
	}

	return targetCompetition, nil
}

func intersection(s1, s2 []string) (inter []string) {
	hash := make(map[string]bool)
	for _, e := range s1 {
		hash[e] = true
	}
	for _, e := range s2 {
		// If elements present in the hashmap then append intersection list.
		if hash[e] {
			inter = append(inter, e)
		}
	}
	//Remove dups from slice.
	inter = removeDups(inter)
	return
}

//Remove dups from slice.
func removeDups(elements []string) (nodups []string) {
	encountered := make(map[string]bool)
	for _, element := range elements {
		if !encountered[element] {
			nodups = append(nodups, element)
			encountered[element] = true
		}
	}
	return
}

func (m MemoryRepository) ListEvents(_ context.Context, first int, page int, sportKey string, categoryKey string, competitionKey string) ([]event.Event, error) {
	if len(m.events) == 0 {
		return nil, ErrEventNotFound
	}
	if sportKey == "" && categoryKey == "" && competitionKey == "" {
		if len(m.events) < (page-1)*first {
			return nil, ErrOutOfRange
		}

		var target []event.Event
		if len(m.events) <= (page)*first {
			target = m.events[(page-1)*first:]
		} else {
			target = m.events[(page-1)*first : (page)*first]
		}

		return target, nil
	}

	var targets []string
	if sportKey != "" {
		v, ok := m.sportsEvents[sportKey]
		if !ok {
			return nil, ErrEventNotFound
		}

		if len(v) < (page-1)*first {
			return nil, ErrOutOfRange
		}
		targets = v
	}

	if categoryKey != "" {
		v, ok := m.categoriesEvents[categoryKey]
		if !ok {
			return nil, ErrEventNotFound
		}

		if len(targets) == 0 {
			targets = v
		} else {
			targets = intersection(targets, v)
			if len(targets) == 0 {
				return nil, ErrEventNotFound
			}
		}
	}

	if competitionKey != "" {
		v, ok := m.competitionsEvents[competitionKey]
		if !ok {
			return nil, ErrEventNotFound
		}

		if len(targets) == 0 {
			targets = v
		} else {
			targets = intersection(targets, v)
			if len(targets) == 0 {
				return nil, ErrEventNotFound
			}
		}
	}

	if len(targets) < (page-1)*first {
		return nil, ErrOutOfRange
	}

	if len(targets) <= (page)*first {
		targets = targets[(page-1)*first:]
	} else {
		targets = targets[(page-1)*first : (page)*first]
	}

	result := make([]event.Event, 0, len(targets))
	for _, target := range targets {
		for _, v := range m.events {
			if target == v.Key() {
				result = append(result, v)
			}
		}
	}
	return result, nil
}

func (m MemoryRepository) GetEvent(_ context.Context, eventKey string) (event.Event, error) {
	var found bool
	var targetEvent event.Event
	for _, v := range m.events {
		if v.Key() == eventKey {
			found = true
			targetEvent = v
			break
		}
	}
	if !found {
		return event.Event{}, ErrEventNotFound
	}

	return targetEvent, nil
}

func (m MemoryRepository) ListEventsCutOffSoon(_ context.Context) ([]event.Event, error) {
	if len(m.events) == 0 {
		return nil, ErrEventNotFound
	}

	var result []event.Event
	for _, event := range m.events {
		if event.Active() && event.CutOffTime().Before(time.Now().Add(5*time.Minute)) {
			result = append(result, event)
		}
	}

	return result, nil
}
