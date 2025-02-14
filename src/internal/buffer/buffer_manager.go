package buffer

import (
	"fmt"
	"sync"
)

type BufferManager struct {
	logDataManager *DataLogsDictionaryManager
	syncMutex      sync.Mutex
	pageDictionary PageDictionary
	pages          []Page
	matches        []Match
}

func NewBufferManager() (BufferManager, error) {
	logDataManager, err := NewDataLogsDictionaryManager()
	if err != nil {
		return BufferManager{}, err
	}

	return BufferManager{
		logDataManager: &logDataManager,
	}, nil
}

func (buffer *BufferManager) AppendToLogsDictionary(correlationId, parsedExpressionType,
	parsedExpression string) error {
	return buffer.logDataManager.Append(NewDataLog(correlationId, parsedExpressionType, parsedExpression))
}

func (buffer *BufferManager) LoadFromLogsDictionary() ([]DataLog, error) {
	if !buffer.logDataManager.IsLogsDictionaryFileExists {
		return []DataLog{}, nil
	}

	result, err := buffer.logDataManager.Read()
	if err != nil {
		return nil, nil
	}
	return result, nil
}

func (buffer *BufferManager) AddPlayer(entityId string, value int16, matchRef *Match) error {
	// move it to processor?
	if _, err := buffer.getPageAdresses(entityId); err == nil {
		return fmt.Errorf("player already exists: %s", entityId)
	}

	pageAdress := buffer.addPage(entityId, value, matchRef)

	buffer.syncMutex.Lock()

	err := buffer.addPageToDictionary(entityId, pageAdress)
	if err != nil {
		return err
	}

	buffer.syncMutex.Unlock()
	return nil
}

func (buffer *BufferManager) UpdatePlayersElo(entityKeys []string, eloLambda int16, matchRef *Match) error {
	for i := range entityKeys {
		err := buffer.UpdatePlayer(entityKeys[i], eloLambda, matchRef)
		if err != nil {
			return err
		}
	}
	return nil
}

func (buffer *BufferManager) UpdatePlayer(entityId string, value int16, matchRef *Match) error {
	isAddedToExistingPage, errTryToAppend := buffer.tryAppendToPage(entityId, value, matchRef)
	if errTryToAppend != nil {
		return errTryToAppend
	}

	if !isAddedToExistingPage {
		pageAdress := buffer.addPage(entityId, value, matchRef)

		buffer.syncMutex.Lock()

		errAppendPageToExistingKey := buffer.appendPageToExistingId(entityId, pageAdress)
		if errAppendPageToExistingKey != nil {
			return errAppendPageToExistingKey
		}

		buffer.syncMutex.Unlock()
	}

	return nil
}

func (buffer *BufferManager) GetPages(entityId string) ([]Page, error) {
	pageAdresses, err := buffer.getPageAdresses(entityId)
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
	entityIdsResult := []string{}
	pagesCollectionResult := [][]Page{}

	for i := range buffer.pageDictionary.pagesData {

		entityIdsResult = append(entityIdsResult, buffer.pageDictionary.pagesData[i].entityId)
		pages := []Page{}
		for j := range buffer.pageDictionary.pagesData[i].pageAdresses {
			pages = append(pages, *buffer.pageDictionary.pagesData[i].pageAdresses[j])
		}
		pagesCollectionResult = append(pagesCollectionResult, pages)
	}

	return entityIdsResult, pagesCollectionResult, nil
}

func (buffer *BufferManager) AddMatchEvent(teamOneKeys []string, teamTwoKeys []string, matchResult int8) *Match {
	newMatch := newMatch(teamOneKeys, teamTwoKeys, matchResult)

	buffer.matches = append(buffer.matches, newMatch)
	return &buffer.matches[len(buffer.matches)-1]
}

func (buffer *BufferManager) addPage(entityId string, value int16, matchRef *Match) *Page {
	buffer.pages = append(buffer.pages, newPage(entityId, value, matchRef))
	return &buffer.pages[len(buffer.pages)-1]
}

func (buffer *BufferManager) tryAppendToPage(entityId string, value int16, matchRef *Match) (bool, error) {
	pageAdresses, err := buffer.getPageAdresses(entityId)
	if err != nil {
		return false, err
	}

	for i := 0; i < len(pageAdresses); i++ {
		if !pageAdresses[i].isPageFull() {
			pageAdresses[i].append(value, matchRef)
			return true, nil
		}
	}
	return false, nil
}
