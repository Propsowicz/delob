package model

import (
	buffer "delob/internal/buffer"
)

type Player struct {
	Key     string
	Elo     int16
	records []int16
}

func NewPlayer(key string, pages []buffer.Page) Player {
	records := []int16{}
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

			records = append(
				records,
				pages[i].Body[j].Value,
			)
		}
	}

	player := Player{
		Key:     key,
		Elo:     elo,
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
