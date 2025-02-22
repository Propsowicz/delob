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

	for i := 0; i < len(pages); i++ {
		for j := 0; j < len(pages[i].Body); j++ {
			if !pages[i].Body[j].IsTransactionStatusSuccessful() {
				continue
			}

			records = append(
				records,
				pages[i].Body[j].Value,
			)
		}
	}

	player := Player{
		Key:     key,
		records: records,
	}
	player.calculateElo()
	return player
}

func (p *Player) calculateElo() {
	var result int16
	for i := range p.records {
		result += p.records[i]
	}

	p.Elo = result
}

func MapPlayerToKeysCollection(players []Player) []string {
	var result []string
	for i := range players {
		result = append(result, players[i].Key)
	}
	return result
}
