package buffer

import (
	"delob/internal/utils"
	"time"
)

type Page struct {
	header Header
	body   Body
}

type Header struct {
	entityId      string
	lastModifed   int64
	isLocked      bool
	lastUsedIndex int16
}

type Body [utils.PAGE_SIZE]Record

type Method int8

const (
	add      Method = iota
	subtract Method = iota
)

type Record struct {
	method Method
	value  float32
}

func (buffer *BufferManager) addPage(entityId string) *Page {
	newPage := Page{
		header: Header{
			entityId:      entityId,
			lastModifed:   time.Now().UTC().UnixMilli(),
			isLocked:      true,
			lastUsedIndex: -1,
		},
		body: [utils.PAGE_SIZE]Record{},
	}
	newPage.body[newPage.header.lastUsedIndex+1] = Record{
		method: add,
		value:  float32(utils.INITIAL_ELO),
	}
	newPage.header.isLocked = false
	newPage.header.lastUsedIndex = newPage.header.lastUsedIndex + 1

	buffer.pages = append(buffer.pages, newPage)
	return &buffer.pages[len(buffer.pages)-1]
}
