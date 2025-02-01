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
	Unknown  Method = iota
	Add      Method = iota
	Subtract Method = iota
)

type Record struct {
	Method Method
	Value  int16
}

func (page *Page) isPageFull() bool {
	return page.Header.lastUsedIndex+1 == utils.PAGE_SIZE
}

func (page *Page) append(record Record) {
	page.Body[page.Header.lastUsedIndex+1] = record
	page.Header.lastUsedIndex++
}

func (buffer *BufferManager) addPage(entityId string, record Record) *Page {
	newPage := Page{
		Header: Header{
			entityId:      entityId,
			lastModifed:   time.Now().UTC().UnixMilli(),
			isLocked:      true,
			lastUsedIndex: -1,
		},
		Body: [utils.PAGE_SIZE]Record{},
	}
	newPage.Body[newPage.Header.lastUsedIndex+1] = record
	newPage.Header.isLocked = false
	newPage.Header.lastUsedIndex = newPage.Header.lastUsedIndex + 1

	buffer.pages = append(buffer.pages, newPage)
	return &buffer.pages[len(buffer.pages)-1]
}

func (buffer *BufferManager) tryAppendToPage(entityId string, record Record) (bool, error) {
	pageAdresses, err := buffer.getPageAdresses(entityId)
	if err != nil {
		return false, err
	}

	for i := 0; i < len(pageAdresses); i++ {
		if !pageAdresses[i].isPageFull() {
			pageAdresses[i].append(record)
			return true, nil
		}
	}
	return false, nil
}
