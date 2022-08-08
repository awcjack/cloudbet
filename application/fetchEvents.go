package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/awcjack/cloudbet/application/query"
	"github.com/awcjack/cloudbet/domain/category"
	"github.com/awcjack/cloudbet/domain/competition"
	"github.com/awcjack/cloudbet/domain/event"
	"github.com/awcjack/cloudbet/domain/sport"
)

type Error struct {
	// Error status
	Status string `json:"status"`
	// Additional Error Details, if applicable for a given error
	Error string `json:"error"`
}

type Sport struct {
	// name of this Sport
	//
	// example: Tennis
	Name string `json:"name"`
	// slug for this Sport
	//
	// example: tennis
	Key string `json:"key"`
	// number of competitions associated with this Sport, 0 indicates inactive Sport
	//
	// example: 2
	CompetitionCount int `json:"competitionCount"`
	// number of events associated with this Sport, 0 indicates inactive Sport
	//
	// example: 4
	EventCount int `json:"eventCount"`
}

type Sports struct {
	// list of all sports offerred
	Sports []Sport `json:"sports"`
}

var (
	ErrMissingToken          = errors.New("missing access token")
	ErrMissingSportKey       = errors.New("missing Sport key")
	ErrMissingCompetitionKey = errors.New("missing competition key")
	ErrNoSports              = errors.New("no sport found")
)

func fetchAllSports(apiKey string) (*Sports, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://sports-api-stg.cloudbet.com/pub/v2/odds/sports", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var response Sports
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	} else {
		var response Error
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		if len(response.Error) > 0 {
			return nil, fmt.Errorf("[%s] %s", response.Status, response.Error)
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

type CompetitionForSport struct {
	// name of this Competition
	//
	// example: "French Open, Men Singles"
	Name string `json:"name"`
	// slug for this Competition. Composed of <sport-key>-<category-key>-<competition-key> as shown in the example value.
	//
	// example: "tennis-atp-french-open-men-singles"
	Key string `json:"key"`
	// number of events associated with this Competition, 0 events indicates inactive Competition
	//
	// example: 2
	EventCount int `json:"eventCount"`
}

type Category struct {
	// name of this Category
	//
	// example: "ATP"
	Name string `json:"name"`
	// slug for this Category
	//
	// example: "atp"
	Key string `json:"key"`
	// list of all competitions associated with this Category
	Competitions []CompetitionForSport `json:"competitions"`
}

type SportWithCategory struct {
	// name of this Sport
	//
	// example: Tennis
	Name string `json:"name"`
	// slug for this Sport
	//
	// example: tennis
	Key string `json:"key"`
	// list of all categories associated with this Sport
	Categories []Category `json:"categories"`
}

func fetchAllCompetitionsUnderSport(apiKey string, sportKey string) (*SportWithCategory, error) {
	if len(sportKey) == 0 {
		return nil, ErrMissingSportKey
	}
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://sports-api-stg.cloudbet.com/pub/v2/odds/sports/%s", sportKey), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var response SportWithCategory
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	} else {
		var response Error
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		if len(response.Error) > 0 {
			return nil, fmt.Errorf("[%s] %s", response.Status, response.Error)
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

type Identifier struct {
	// name of this Identifier
	Name string `json:"name"`
	// slug for this Identifier
	Key string `json:"key"`
}

type CompetitionWithCategory struct {
	// name of this Competition
	Name string `json:"name"`
	// slug for this Competition
	Key string `json:"key"`
	// category associated with this Competition
	Category *Identifier `json:"category"`
}

type TeamIdentifier struct {
	// name of this Identifier
	Name string `json:"name"`
	// slug for this Identifier
	Key string `json:"key"`
	// abbreviation for this team's name
	Abbreviation string `json:"abbreviation"`
	// team country code
	Nationality string `json:"nationality"`
}

type Market struct {
	// mapping between submarket key and all associated submarkets for this Market
	Submarkets map[string]*Submarket `json:"submarkets"`
}

type Submarket struct {
	// sequential update number for this Submarket
	Sequence string `json:"sequence"`
	// list of all associated selections for this Submarket
	Selections []Selection `json:"selections"`
}

type Selection struct {
	// outcome of this Selection
	Outcome string `json:"outcome"`
	// parameters to be sent by the client during bet placement on this selection, such as handicap, period etc.
	Params string `json:"params"`
	// price at which bets can be placed on this Selection
	Price float64 `json:"price"`
	// maximum stake in EUR which can be placed in bets on this Selection; market liability = selection max stake * (price - 1); minimum stake is 0.01 EUR for all markets
	MaxStake float64 `json:"maxStake"`
	// probability of this Selection's outcome
	Probability float64 `json:"probability"`
	// current status of this Selection
	Status string `json:"status"`
	// side of this Selection (back/lay)
	Side string `json:"side"`
}

type EventMetadata struct {
	// opinion is an answer to question "how players bet"
	Opinion []OutcomeProbability `json:"opinion"`
}

type OutcomeProbability struct {
	// outcome
	Outcome string `json:"outcome"`
	// probability
	Probability float32 `json:"probability"`
	// market key used to build opinion
	MarketKey string `json:"marketKey"`
	// market parameters used to build opinion, such as handicap, period etc.
	Params string `json:"params"`
}

type Event struct {
	// sequential update number for this Event
	Sequence string `json:"sequence"`
	// unique ID for this Event
	Id int `json:"id"`
	// sport associated with this Event
	Sport *Identifier `json:"sport"`
	// competition associated with this Event
	Competition *CompetitionWithCategory `json:"competition"`
	// key-name tuple for the home team competitor of this Event
	Home *TeamIdentifier `json:"home"`
	// key-name tuple for the away team competitor of this Event
	Away *TeamIdentifier `json:"away"`
	// current status of this Event
	Status string `json:"status"`
	// mapping between market key and all associated markets for this Event
	Markets map[string]*Market `json:"markets"`
	// name of this Event
	Name string `json:"name"`
	// slug for this Event
	Key string `json:"key"`
	// event cutoff time in string format "2006-01-02T15:04:05Z07:00" (RFC3339)
	CutoffTime string `json:"cutoffTime"`
	// metadata
	Metadata *EventMetadata `json:"metadata"`
}

type Competition struct {
	// name of this Competition
	//
	// example: "French Open, Men Singles"
	Name string `json:"name"`
	// slug for this Competition. Composed of <sport-key>-<category-key>-<competition-key> as shown in the example value.
	//
	// example: "tennis-atp-french-open-men-singles"
	Key string `json:"key"`
	// sport associated with this Competition
	//
	// example: {"name":"Tennis","key":"tennis"}
	Sport Identifier `json:"sport"`
	// list of all events associated with this competition
	Events []Event `json:"events"`
	// category associated with this Competition
	//
	// example: {"name":"ATP","key":"atp"}
	Category Identifier `json:"category"`
}

func fetchAllEventsUnderCompetition(apiKey string, competitionKey string) (*Competition, error) {
	if len(competitionKey) == 0 {
		return nil, ErrMissingCompetitionKey
	}
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://sports-api-stg.cloudbet.com/pub/v2/odds/competitions/%s", competitionKey), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var response Competition
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println(string(body))
			return nil, err
		}
		return &response, nil
	} else {
		var response Error
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		if len(response.Error) > 0 {
			return nil, fmt.Errorf("[%s] %s", response.Status, response.Error)
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

type logger interface {
	Panicf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Warningf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

type CloudbetHandler struct {
	eventRepo event.Repository
	logger    logger
	apiKey    string
}

func NewCloudbetHander(eventRepo event.Repository, logger logger, apiKey string) CloudbetHandler {
	return CloudbetHandler{
		eventRepo: eventRepo,
		logger:    logger,
		apiKey:    apiKey,
	}
}

func (h CloudbetHandler) StoreAllEvents(ctx context.Context) error {
	if len(h.apiKey) == 0 {
		h.logger.Panicf("Missing cloudbet access token")
	}
	allSports, err := fetchAllSports(h.apiKey)
	if err != nil {
		h.logger.Errorf("Fetch All Sports Error %s", err)
	}

	if len(allSports.Sports) == 0 {
		h.logger.Warningf("No sports found")
		return ErrNoSports
	}

	h.logger.Debugf("Fetched All Sports: %v", allSports.Sports)

	for _, sport := range allSports.Sports {
		if sport.CompetitionCount != 0 && sport.EventCount != 0 {
			// only check active sport
			if len(sport.Key) == 0 {
				h.logger.Errorf("Missing Sport key for querying competition")
				return ErrMissingSportKey
			}

			categorizedSport, err := fetchAllCompetitionsUnderSport(h.apiKey, sport.Key)
			if err != nil {
				h.logger.Errorf("Fetch competitions under sport %s Error %s", sport.Name, err)
				continue
			}

			h.logger.Debugf("Fetched All categories of sport %s: %v", sport.Name, categorizedSport.Categories)
			if len(categorizedSport.Categories) == 0 {
				h.logger.Warningf("No categories found under sport %s", sport.Key)
				continue
			}

			sportIdentity, err := event.NewIdentifier(sport.Name, sport.Key)
			if err != nil {
				h.logger.Errorf("Cannot create sport identity %s", sport.Name)
				continue
			}
			for _, category := range categorizedSport.Categories {
				categoryIdentity, err := event.NewIdentifier(category.Name, category.Key)
				if err != nil {
					h.logger.Errorf("Cannot create category identity %s", category.Name)
					continue
				}
				for _, competition := range category.Competitions {
					if len(competition.Key) == 0 {
						h.logger.Errorf("Missing competition key for querying event")
						continue
					}
					if competition.EventCount == 0 {
						h.logger.Warningf("Event not found under competition %s", competition.Name)
						continue
					}

					competitionIdentity, err := event.NewIdentifier(competition.Name, competition.Key)
					if err != nil {
						h.logger.Errorf("Cannot create competition identity %s", competition.Name)
						continue
					}
					competitionInfo, err := fetchAllEventsUnderCompetition(h.apiKey, competition.Key)
					if err != nil {
						h.logger.Errorf("Fetch event under sport %s competition %s Error %s", sport.Name, competition.Name, err)
						continue
					}
					for _, competitionEvent := range competitionInfo.Events {
						var homeIdentity event.TeamIdentifier
						if competitionEvent.Home != nil {
							homeIdentity = event.NewTeamIdentifier(competitionEvent.Home.Name, competitionEvent.Home.Key, competitionEvent.Home.Abbreviation, competitionEvent.Home.Nationality)
						} else {
							homeIdentity = event.NewTeamIdentifier("", "", "", "")
						}
						var awayIdentity event.TeamIdentifier
						if competitionEvent.Home != nil {
							awayIdentity = event.NewTeamIdentifier(competitionEvent.Away.Name, competitionEvent.Away.Key, competitionEvent.Away.Abbreviation, competitionEvent.Away.Nationality)
						} else {
							awayIdentity = event.NewTeamIdentifier("", "", "", "")
						}

						var active bool
						if competitionEvent.Status == "TRADING" || competitionEvent.Status == "TRADING_LIVE" {
							active = true
						}

						marketValue := make(map[string]event.Market)
						for key, market := range competitionEvent.Markets {
							subMarketValue := make(map[string][]event.Selection)
							for subMarketKey, subMarket := range market.Submarkets {
								selections := make([]event.Selection, len(subMarket.Selections))
								for index, selection := range subMarket.Selections {
									selections[index] = event.NewSelection(selection.Outcome, selection.Params, selection.Price, selection.MaxStake, selection.Probability, selection.Status, selection.Side)
								}
								subMarketValue[subMarketKey] = selections
							}
							marketValue[key] = event.NewMarket(subMarketValue)
						}

						cutOffTime, err := time.Parse(time.RFC3339, competitionEvent.CutoffTime)
						if err != nil {
							h.logger.Errorf("Parse Cut off time error %s", err)
							continue
						}

						e, err := event.NewEvent(&sportIdentity, &competitionIdentity, &categoryIdentity, homeIdentity, awayIdentity, active, marketValue, competitionEvent.Name, competitionEvent.Key, cutOffTime)
						if err != nil {
							h.logger.Errorf("Parse Cut off time error %s", err)
							continue
						}

						err = h.eventRepo.Save(ctx, *e)
						if err != nil {
							h.logger.Errorf("Store data error %s", err)
						}
					}
				}

			}
		} else {
			h.logger.Debugf("Skip Inactive sport %s", sport.Name)
		}
	}

	h.logger.Infof("Finished Fetch Events from Cloudbet and store to DB")
	return nil
}

func fetchEventInfo(apiKey string, eventKey string) (*Event, error) {
	if eventKey == "" {
		return nil, ErrMissingCompetitionKey
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://sports-api-stg.cloudbet.com/pub/v2/odds/events/%s", eventKey), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var response Event
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println(string(body))
			return nil, err
		}
		return &response, nil
	} else {
		var response Error
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		if len(response.Error) > 0 {
			return nil, fmt.Errorf("[%s] %s", response.Status, response.Error)
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

func (h CloudbetHandler) CheckEventsCloseToCutOff(ctx context.Context) {
	if len(h.apiKey) == 0 {
		h.logger.Panicf("Missing cloudbet access token")
	}
	events, err := h.eventRepo.ListEventsCutOffSoon(ctx)
	if err != nil {
		h.logger.Errorf(err.Error())
	}

	for _, event := range events {
		latestEvent, err := fetchEventInfo(h.apiKey, event.Key())
		if err != nil {
			h.logger.Errorf(err.Error())
			continue
		}

		if latestEvent.Status != "TRADING" && latestEvent.Status != "TRADING_LIVE" {
			event.Inactivate()
		}
	}
}

type Queries struct {
	ListSports       *query.ListSportsHandler
	GetSport         *query.GetSportHandler
	ListCategories   *query.ListCategoriesHandler
	GetCategory      *query.GetCategoryHandler
	ListCompetitions *query.ListCompetitionsHandler
	GetCompetition   *query.GetCompetitionHandler
	ListEvents       *query.ListEventsHandler
	GetEvent         *query.GetEventHandler
}

type Application struct {
	Query Queries
}

func NewApplication(sportRepo sport.Repository, categoryRepo category.Repository, competitionRepo competition.Repository, eventRepo event.Repository, logger logger) *Application {
	listSportsHandler := query.NewListSportHandler(sportRepo, logger)
	getSportHandler := query.NewGetSportHandler(sportRepo, logger)
	listCategoriesHandler := query.NewListCategoriesHandler(categoryRepo, logger)
	getCategoryHandler := query.NewGetCategoryHandler(categoryRepo, logger)
	listCompetitionsHandler := query.NewListCompetitionHandler(competitionRepo, logger)
	getCompetitionHandler := query.NewGetCompetitionHandler(competitionRepo, logger)
	listEventsHandler := query.NewListEventsHandler(eventRepo, logger)
	getEventHandler := query.NewGetEventHandler(eventRepo, logger)

	return &Application{
		Query: Queries{
			ListSports:       listSportsHandler,
			GetSport:         getSportHandler,
			ListCategories:   listCategoriesHandler,
			GetCategory:      getCategoryHandler,
			ListCompetitions: listCompetitionsHandler,
			GetCompetition:   getCompetitionHandler,
			ListEvents:       listEventsHandler,
			GetEvent:         getEventHandler,
		},
	}
}
