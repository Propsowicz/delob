package processor

import (
	buffer "delob/internal/buffer"
)

type Player struct {
	Id      string
	Elo     int16
	records []int16
}

func newPlayer(id string, pages []buffer.Page) Player {
	primitiveRecords := []int16{}

	for i := 0; i < len(pages); i++ {
		for j := 0; j < len(pages[i].Body); j++ {
			primitiveRecords = append(
				primitiveRecords,
				pages[i].Body[j].Value,
			)
		}
	}

	player := Player{
		Id:      id,
		records: primitiveRecords,
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
