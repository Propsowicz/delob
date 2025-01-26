package buffer

import (
	"delob/internal/utils"
	"time"
)

type Page struct {
	Header Header
	Body   Body
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
	Add      Method = iota
	Subtract Method = iota
)

type Record struct {
	Method Method
	Value  int16
}

func (buffer *BufferManager) addPage(entityId string) *Page {
	newPage := Page{
		Header: Header{
			entityId:      entityId,
			lastModifed:   time.Now().UTC().UnixMilli(),
			isLocked:      true,
			lastUsedIndex: -1,
		},
		Body: [utils.PAGE_SIZE]Record{},
	}
	newPage.Body[newPage.Header.lastUsedIndex+1] = Record{
		Method: Add,
		// TODO move it out of here
		Value: utils.INITIAL_ELO,
	}
	newPage.Header.isLocked = false
	newPage.Header.lastUsedIndex = newPage.Header.lastUsedIndex + 1

	buffer.pages = append(buffer.pages, newPage)
	return &buffer.pages[len(buffer.pages)-1]
}
