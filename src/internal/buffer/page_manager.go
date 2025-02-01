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
	isCached      bool
	cachedValue   int16
	entityId      string
	lastModifed   int64
	isLocked      bool
	lastUsedIndex int16
}

type Body [utils.PAGE_SIZE]Record

type Record struct {
	IsUsed       bool
	AddTimestamp int64
	Value        int16
}

func (page *Page) isPageFull() bool {
	return page.Header.lastUsedIndex+1 == utils.PAGE_SIZE
}

func (page *Page) append(value int16) {
	page.Body[page.Header.lastUsedIndex+1] = newRecord(value)
	page.Header.lastUsedIndex++
}

func newPage(entityId string, value int16) Page {
	newPage := Page{
		Header: Header{
			entityId:      entityId,
			lastModifed:   time.Now().UTC().UnixMilli(),
			isLocked:      false,
			lastUsedIndex: 0,
		},
		Body: [utils.PAGE_SIZE]Record{newRecord(value)},
	}

	return newPage
}

func newRecord(value int16) Record {
	return Record{
		IsUsed:       true,
		AddTimestamp: time.Now().UnixMilli(),
		Value:        value,
	}
}
