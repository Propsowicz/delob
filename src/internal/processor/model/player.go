package model

import (
	buffer "delob/internal/buffer"
	"delob/internal/parser"
	"delob/internal/shared"
	"delob/internal/utils"
	"time"
)

type Player struct {
	Key     string  `json:"Key"`
	Elo     int16   `json:"Elo,omitempty"`
	Events  []Event `json:"Events,omitempty"`
	records []int16
}

type Event struct {
	Change      int16 `json:"Change,omitempty"`
	DateTime    time.Time
	TeamOne     []string `json:"TeamOne,omitempty"`
	TeamTwo     []string `json:"TeamTwo,omitempty"`
	MatchResult string   `json:"MatchResult,omitempty"`
}

func NewPlayer(key string, pages []buffer.Page, queryComponents []parser.SelectQueryComponent) Player {
	records := []int16{}
	eventsDto := []Event{}
	var calculatedElo int16
	var eloDto int16
	var keyDto string

	selectEvents := utils.Contains(queryComponents, parser.EventsComponent)
	selectMatches := utils.Contains(queryComponents, parser.MatchesComponent)

	for i := 0; i < len(pages); i++ {

		if pages[i].Header.IsCached {
			calculatedElo += pages[i].Header.CachedValue
		}

		for j := 0; j < len(pages[i].Body); j++ {
			if !pages[i].Body[j].IsTransactionStatusSuccessful() {
				continue
			}

			if !pages[i].Header.IsCached {
				calculatedElo += pages[i].Body[j].Value
			}

			if selectEvents || selectMatches {
				eventsDto = append(eventsDto, parseEvents(selectEvents, selectMatches, *pages[i].Body[j]))
			}
		}
	}

	if utils.Contains(queryComponents, parser.KeyComponent) {
		keyDto = key
	}
	if utils.Contains(queryComponents, parser.EloComponent) {
		eloDto = calculatedElo
	}

	player := Player{
		Key:     keyDto,
		Elo:     eloDto,
		Events:  eventsDto,
		records: records,
	}
	return player
}

func parseEvents(selectEvents, selectMatches bool, body buffer.Record) Event {
	eventDto := Event{}

	eventDto.DateTime = time.UnixMilli(body.AddTimestamp)

	if selectEvents {
		eventDto.Change = body.Value
	}

	if selectMatches && body.MatchRef != nil {
		eventDto.TeamOne = body.MatchRef.TeamOneKeys
		eventDto.TeamTwo = body.MatchRef.TeamTwoKeys
		eventDto.MatchResult = shared.MapToString(shared.MatchResult(body.MatchRef.MatchResult))
	}

	return eventDto
}

func MapPlayerToKeysCollection(players []Player) []string {
	var result []string
	for i := range players {
		result = append(result, players[i].Key)
	}
	return result
}
