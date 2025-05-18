package model

import (
	buffer "delob/internal/buffer"
	"time"
)

type Player struct {
	Key     string  `json:"Key"`
	Elo     int16   `json:"Elo"`
	Stats   []Stats `json:"Stats,omitempty"`
	records []int16
}

type Stats struct {
	Change   int16
	DateTime time.Time
}

// type Player struct {
//     ID     string  `json:"id"`
//     Name   string  `json:"name,omitempty"`
//     Score  *int    `json:"score,omitempty"`
//     Active bool    `json:"active,omitempty"`
// }

func NewPlayer(key string, pages []buffer.Page) Player {
	records := []int16{}
	stats := []Stats{}
	var elo int16

	for i := 0; i < len(pages); i++ {

		if pages[i].Header.IsCached {
			elo += pages[i].Header.CachedValue
		}

		for j := 0; j < len(pages[i].Body); j++ {
			if !pages[i].Body[j].IsTransactionStatusSuccessful() {
				continue
			}

			if !pages[i].Header.IsCached {
				elo += pages[i].Body[j].Value
			}
			// stats = append(stats,
			// 	Stats{
			// 		Change:   pages[i].Body[j].Value,
			// 		DateTime: time.UnixMilli(pages[i].Body[j].AddTimestamp),
			// 	},
			// )
			records = append(
				records,
				pages[i].Body[j].Value,
			)
		}
	}

	player := Player{
		Key:     key,
		Elo:     elo,
		Stats:   stats,
		records: records,
	}
	return player
}

func MapPlayerToKeysCollection(players []Player) []string {
	var result []string
	for i := range players {
		result = append(result, players[i].Key)
	}
	return result
}
