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
	key           string
	lastModifed   int64
	isLocked      bool
	lastUsedIndex int16
}

type Body [utils.PAGE_SIZE]*Record

type Record struct {
	transactionStatus transactionStatus
	AddTimestamp      int64
	Value             int16
	MatchRefKey       string
	MatchRef          *Match
}

func (r *Record) IsTransactionStatusSuccessful() bool {
	return r.transactionStatus == success
}

func (page *Page) isPageFull() bool {
	return page.Header.lastUsedIndex+1 == utils.PAGE_SIZE
}

func (page *Page) append(value int16, matchRef *Match, transaction *Transaction) {
	overrideEmptyRecord(page.Body[page.Header.lastUsedIndex+1], value, matchRef, transaction)
	page.Header.lastUsedIndex++
}

func newPage(key string, value int16, matchRef *Match, transaction *Transaction) *Page {
	newPage := Page{
		Header: Header{
			key:           key,
			lastModifed:   utils.Timestamp(),
			isLocked:      false,
			lastUsedIndex: 0,
		},
		Body: [utils.PAGE_SIZE]*Record{},
	}

	for i := range len(newPage.Body) {
		if i == 0 {
			newPage.Body[i] = newRecord(value, matchRef, transaction)
			continue
		}
		newPage.Body[i] = &Record{}
	}

	return &newPage
}

func newRecord(value int16, matchRef *Match, transaction *Transaction) *Record {
	matchRefKey := "init"
	if matchRef != nil {
		matchRefKey = matchRef.Key
	}

	record := &Record{
		transactionStatus: inProgress,
		AddTimestamp:      utils.Timestamp(),
		Value:             value,
		MatchRef:          matchRef,
		MatchRefKey:       matchRefKey,
	}

	transaction.AddRecordPointer(record)

	return record
}

func overrideEmptyRecord(r *Record, value int16, matchRef *Match, transaction *Transaction) *Record {
	matchRefKey := "init"
	if matchRef != nil {
		matchRefKey = matchRef.Key
	}

	r.transactionStatus = inProgress
	r.AddTimestamp = utils.Timestamp()
	r.Value = value
	r.MatchRef = matchRef
	r.MatchRefKey = matchRefKey

	transaction.AddRecordPointer(r)
	return r
}

func (buffer *BufferManager) addPage(key string, value int16, matchRef *Match, transaction *Transaction) *Page {
	buffer.pages = append(buffer.pages, newPage(key, value, matchRef, transaction))
	return buffer.pages[len(buffer.pages)-1]
}

func (buffer *BufferManager) tryAppendToPage(key string, value int16, matchRef *Match, transaction *Transaction) (bool, error) {
	pageAdresses, err := buffer.getPageAdresses(key, isSuccess)
	if err != nil {
		return false, err
	}

	for i := 0; i < len(pageAdresses); i++ {
		if !pageAdresses[i].isPageFull() {
			pageAdresses[i].append(value, matchRef, transaction)
			return true, nil
		}
	}
	return false, nil
}

func (buffer *BufferManager) GetPages(key string) ([]Page, error) {
	pageAdresses, err := buffer.getPageAdresses(key, isSuccess)
	if err != nil {
		return []Page{}, err
	}
	result := []Page{}

	for i := range pageAdresses {
		result = append(result, *pageAdresses[i])
	}
	return result, nil
}

func (buffer *BufferManager) GetAllPages() ([]string, [][]Page, error) {
	keysResult := []string{}
	pagesCollectionResult := [][]Page{}

	for i := range buffer.pageDictionary.pagesData {
		if !isSuccess(buffer.pageDictionary.pagesData[i].transactionStatus) {
			continue
		}
		keysResult = append(keysResult, buffer.pageDictionary.pagesData[i].key)
		pages := []Page{}
		for j := range buffer.pageDictionary.pagesData[i].pageAdresses {
			pages = append(pages, *buffer.pageDictionary.pagesData[i].pageAdresses[j])
		}
		pagesCollectionResult = append(pagesCollectionResult, pages)
	}

	return keysResult, pagesCollectionResult, nil
}
