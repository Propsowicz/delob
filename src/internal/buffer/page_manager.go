package buffer

import (
	"delob/internal/utils"
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
	MatchRefKey  string
	MatchRef     *Match
}

func (page *Page) isPageFull() bool {
	return page.Header.lastUsedIndex+1 == utils.PAGE_SIZE
}

func (page *Page) append(value int16, matchRef *Match) {
	page.Body[page.Header.lastUsedIndex+1] = newRecord(value, matchRef)
	page.Header.lastUsedIndex++
}

func newPage(entityId string, value int16, matchRef *Match) Page {
	newPage := Page{
		Header: Header{
			entityId:      entityId,
			lastModifed:   utils.Timestamp(),
			isLocked:      false,
			lastUsedIndex: 0,
		},
		Body: [utils.PAGE_SIZE]Record{newRecord(value, matchRef)},
	}

	return newPage
}

func newRecord(value int16, matchRef *Match) Record {
	matchRefKey := "init"
	if matchRef != nil {
		matchRefKey = matchRef.Key
	}

	return Record{
		IsUsed:       true,
		AddTimestamp: utils.Timestamp(),
		Value:        value,
		MatchRef:     matchRef,
		MatchRefKey:  matchRefKey,
	}
}
