//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.11.0 --config=codegen.config.yaml ../docs/openapi.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.11.0 --config=typegen.config.yaml ../docs/openapi.yaml

package interfaces

import (
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/awcjack/cloudbet/application"
	"github.com/gin-gonic/gin"
)

var (
	ErrSizeTooLarge = errors.New("size cannot be smaller than 1 or larger than 50")
	ErrPageTooSmall = errors.New("page cannot be smaller than 1")
	ErrMissingKey   = errors.New("missing key")
)

type HttpServer struct {
	app application.Application
}

func NewHttpServer(app application.Application) *HttpServer {
	return &HttpServer{
		app: app,
	}
}

func (h HttpServer) ListSports(c *gin.Context, params ListSportsParams) {
	if params.First <= 0 || params.First > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrSizeTooLarge.Error()})
		return
	}
	if params.Page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrPageTooSmall.Error()})
		return
	}

	repoData, err := h.app.Query.ListSports.Handle(c, int(params.First), int(params.Page))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := make([]Sport, len(repoData))
	for i, sport := range repoData {
		name := sport.Name()
		liveTime := sport.LiveTime()
		if math.IsNaN(liveTime) {
			liveTime = 0
		}
		result[i] = Sport{
			Key:      sport.Key(),
			Name:     &name,
			LiveTime: &liveTime,
		}
	}
	c.JSON(http.StatusOK, result)
}

func (h HttpServer) GetSport(c *gin.Context, sportKey string) {
	if sportKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrMissingKey.Error()})
		return
	}

	sport, err := h.app.Query.GetSport.Handle(c, sportKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := sport.Name()
	liveTime := sport.LiveTime()
	if math.IsNaN(liveTime) {
		liveTime = 0
	}
	c.JSON(http.StatusOK, Sport{
		Key:      sport.Key(),
		Name:     &name,
		LiveTime: &liveTime,
	})
}

func (h HttpServer) ListCategories(c *gin.Context, params ListCategoriesParams) {
	if params.First <= 0 || params.First > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrSizeTooLarge.Error()})
		return
	}
	if params.Page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrPageTooSmall.Error()})
		return
	}

	sportKey := ""
	if params.Sport != nil {
		sportKey = *params.Sport
	}
	repoData, err := h.app.Query.ListCategories.Handle(c, int(params.First), int(params.Page), sportKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := make([]Category, len(repoData))
	for i, category := range repoData {
		name := category.Name()
		result[i] = Category{
			Key:  category.Key(),
			Name: &name,
		}
	}
	c.JSON(http.StatusOK, result)
}

func (h HttpServer) GetCategory(c *gin.Context, categoryKey string) {
	if categoryKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrMissingKey.Error()})
		return
	}

	category, err := h.app.Query.GetCategory.Handle(c, categoryKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := category.Name()
	c.JSON(http.StatusOK, Category{
		Key:  category.Key(),
		Name: &name,
	})
}

func (h HttpServer) ListCompetitions(c *gin.Context, params ListCompetitionsParams) {
	if params.First <= 0 || params.First > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrSizeTooLarge.Error()})
		return
	}
	if params.Page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrPageTooSmall.Error()})
		return
	}

	sportKey := ""
	if params.Sport != nil {
		sportKey = *params.Sport
	}
	repoData, err := h.app.Query.ListCompetitions.Handle(c, int(params.First), int(params.Page), sportKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := make([]Competition, len(repoData))
	for i, competition := range repoData {
		name := competition.Name()
		result[i] = Competition{
			Key:  competition.Key(),
			Name: &name,
		}
	}
	c.JSON(http.StatusOK, result)
}

func (h HttpServer) GetCompetition(c *gin.Context, competitionKey string) {
	if competitionKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrMissingKey.Error()})
		return
	}

	competition, err := h.app.Query.GetCompetition.Handle(c, competitionKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := competition.Name()
	c.JSON(http.StatusOK, Competition{
		Key:  competition.Key(),
		Name: &name,
	})
}

func (h HttpServer) ListEvents(c *gin.Context, params ListEventsParams) {
	if params.First <= 0 || params.First > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrSizeTooLarge.Error()})
		return
	}
	if params.Page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrPageTooSmall.Error()})
		return
	}

	sportKey := ""
	if params.Sport != nil {
		sportKey = *params.Sport
	}
	categoryKey := ""
	if params.Category != nil {
		categoryKey = *params.Category
	}
	competitionKey := ""
	if params.Competition != nil {
		competitionKey = *params.Competition
	}
	repoData, err := h.app.Query.ListEvents.Handle(c, int(params.First), int(params.Page), sportKey, categoryKey, competitionKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := make([]Event, len(repoData))
	for i, event := range repoData {
		name := event.Name()
		sportKey := event.Sport().Key()
		sportName := event.Sport().Name()
		categoryKey := event.Category().Key()
		categoryName := event.Category().Name()
		competitionKey := event.Competition().Key()
		competitionName := event.Competition().Name()
		homeName := event.Home().Name()
		homeAbbreviation := event.Home().Abbreviation()
		homeNationality := event.Home().Nationality()
		awayName := event.Away().Name()
		awayAbbreviation := event.Away().Abbreviation()
		awayNationality := event.Away().Nationality()
		var market Event_Market
		market.AdditionalProperties = make(map[string]Market)
		for k, v := range event.Market() {
			marketProperty := Market{}
			subMarket := make(map[string][]Selection, len(v.Submarkets()))
			for subMarketKey, subMarketVal := range v.Submarkets() {
				subMarket[subMarketKey] = make([]Selection, len(subMarketVal))
				for i, selectionVal := range subMarketVal {
					outcome := selectionVal.Outcome()
					params := selectionVal.Params()
					price := selectionVal.Price()
					maxStake := selectionVal.MaxStake()
					probability := selectionVal.Probability()
					status := selectionVal.Status()
					var statusEnum SelectionStatus
					if status == "SELECTION_DISABLED" {
						statusEnum = SELECTIONDISABLED
					} else {
						statusEnum = SELECTIONENABLED
					}
					side := selectionVal.Side()
					var sideEnum SelectionSide
					if side == "BACK" {
						sideEnum = BACK
					} else {
						sideEnum = LAY
					}

					subMarket[subMarketKey][i] = Selection{
						Outcome:     &outcome,
						Params:      &params,
						Price:       &price,
						MaxStake:    &maxStake,
						Probability: &probability,
						Status:      &statusEnum,
						Side:        &sideEnum,
					}
				}
			}
			marketProperty.Submarkets = &Market_Submarkets{
				AdditionalProperties: subMarket,
			}

			market.AdditionalProperties[k] = marketProperty
		}
		var cutOffTime string
		if event.CutOffTime().IsZero() {
			cutOffTime = ""
		} else {
			cutOffTime = event.CutOffTime().Format(time.RFC3339)
		}
		var startTradingLiveTime string
		if event.StartTradingLiveTime().IsZero() {
			startTradingLiveTime = ""
		} else {
			startTradingLiveTime = event.StartTradingLiveTime().Format(time.RFC3339)
		}
		var inactiveTime string
		if event.InactiveTime().IsZero() {
			inactiveTime = ""
		} else {
			inactiveTime = event.InactiveTime().Format(time.RFC3339)
		}

		result[i] = Event{
			Sport: &struct {
				Key  *string `json:"key,omitempty"`
				Name *string `json:"name,omitempty"`
			}{
				Key:  &sportKey,
				Name: &sportName,
			},
			Category: &struct {
				Key  *string `json:"key,omitempty"`
				Name *string `json:"name,omitempty"`
			}{
				Key:  &categoryKey,
				Name: &categoryName,
			},
			Competition: &struct {
				Key  *string `json:"key,omitempty"`
				Name *string `json:"name,omitempty"`
			}{
				Key:  &competitionKey,
				Name: &competitionName,
			},
			Home: &Team{
				Key:          event.Home().Key(),
				Name:         &homeName,
				Abbreviation: &homeAbbreviation,
				Nationality:  &homeNationality,
			},
			Away: &Team{
				Key:          event.Home().Key(),
				Name:         &awayName,
				Abbreviation: &awayAbbreviation,
				Nationality:  &awayNationality,
			},
			Active:               event.Active(),
			Market:               &market,
			Name:                 &name,
			Key:                  event.Key(),
			CutOffTime:           &cutOffTime,
			StartTradingLiveTime: &startTradingLiveTime,
			InactiveTime:         &inactiveTime,
		}
	}
	c.JSON(http.StatusOK, result)
}

func (h HttpServer) GetEvent(c *gin.Context, eventKey string) {
	if eventKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrMissingKey.Error()})
		return
	}

	event, err := h.app.Query.GetEvent.Handle(c, eventKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := event.Name()
	sportKey := event.Sport().Key()
	sportName := event.Sport().Name()
	categoryKey := event.Category().Key()
	categoryName := event.Category().Name()
	competitionKey := event.Competition().Key()
	competitionName := event.Competition().Name()
	homeName := event.Home().Name()
	homeAbbreviation := event.Home().Abbreviation()
	homeNationality := event.Home().Nationality()
	awayName := event.Away().Name()
	awayAbbreviation := event.Away().Abbreviation()
	awayNationality := event.Away().Nationality()
	var market Event_Market
	market.AdditionalProperties = make(map[string]Market)
	for k, v := range event.Market() {
		marketProperty := Market{}
		subMarket := make(map[string][]Selection, len(v.Submarkets()))
		for subMarketKey, subMarketVal := range v.Submarkets() {
			subMarket[subMarketKey] = make([]Selection, len(subMarketVal))
			for i, selectionVal := range subMarketVal {
				outcome := selectionVal.Outcome()
				params := selectionVal.Params()
				price := selectionVal.Price()
				maxStake := selectionVal.MaxStake()
				probability := selectionVal.Probability()
				status := selectionVal.Status()
				var statusEnum SelectionStatus
				if status == "SELECTION_DISABLED" {
					statusEnum = SELECTIONDISABLED
				} else {
					statusEnum = SELECTIONENABLED
				}
				side := selectionVal.Side()
				var sideEnum SelectionSide
				if side == "BACK" {
					sideEnum = BACK
				} else {
					sideEnum = LAY
				}

				subMarket[subMarketKey][i] = Selection{
					Outcome:     &outcome,
					Params:      &params,
					Price:       &price,
					MaxStake:    &maxStake,
					Probability: &probability,
					Status:      &statusEnum,
					Side:        &sideEnum,
				}
			}
		}
		marketProperty.Submarkets = &Market_Submarkets{
			AdditionalProperties: subMarket,
		}

		market.AdditionalProperties[k] = marketProperty
	}
	var cutOffTime string
	if event.CutOffTime().IsZero() {
		cutOffTime = ""
	} else {
		cutOffTime = event.CutOffTime().Format(time.RFC3339)
	}
	var startTradingLiveTime string
	if event.StartTradingLiveTime().IsZero() {
		startTradingLiveTime = ""
	} else {
		startTradingLiveTime = event.StartTradingLiveTime().Format(time.RFC3339)
	}
	var inactiveTime string
	if event.InactiveTime().IsZero() {
		inactiveTime = ""
	} else {
		inactiveTime = event.InactiveTime().Format(time.RFC3339)
	}
	c.JSON(http.StatusOK, Event{
		Sport: &struct {
			Key  *string `json:"key,omitempty"`
			Name *string `json:"name,omitempty"`
		}{
			Key:  &sportKey,
			Name: &sportName,
		},
		Category: &struct {
			Key  *string `json:"key,omitempty"`
			Name *string `json:"name,omitempty"`
		}{
			Key:  &categoryKey,
			Name: &categoryName,
		},
		Competition: &struct {
			Key  *string `json:"key,omitempty"`
			Name *string `json:"name,omitempty"`
		}{
			Key:  &competitionKey,
			Name: &competitionName,
		},
		Home: &Team{
			Key:          event.Home().Key(),
			Name:         &homeName,
			Abbreviation: &homeAbbreviation,
			Nationality:  &homeNationality,
		},
		Away: &Team{
			Key:          event.Home().Key(),
			Name:         &awayName,
			Abbreviation: &awayAbbreviation,
			Nationality:  &awayNationality,
		},
		Active:               event.Active(),
		Market:               &market,
		Name:                 &name,
		Key:                  event.Key(),
		CutOffTime:           &cutOffTime,
		StartTradingLiveTime: &startTradingLiveTime,
		InactiveTime:         &inactiveTime,
	})
}

func NewHandler(httpServer HttpServer) *gin.Engine {
	router := gin.Default()

	RegisterHandlers(router, httpServer)

	return router
}
