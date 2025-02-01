package buffer

import (
	"fmt"
	"sync"
)

type BufferManager struct {
	syncMutex      sync.Mutex
	pageDictionary PageDictionary
	pages          []Page
}

func NewBufferManager() BufferManager {
	return BufferManager{}
}

func (buffer *BufferManager) AddPlayer(entityId string, value int16) error {

	// move it to processor?
	if _, err := buffer.getPageAdresses(entityId); err == nil {
		return fmt.Errorf("player already exists: %s", entityId)
	}

	pageAdress := buffer.addPage(entityId, value)

	buffer.syncMutex.Lock()

	err := buffer.addPageToDictionary(entityId, pageAdress)
	if err != nil {
		return err
	}

	buffer.syncMutex.Unlock()
	return nil
}

func (buffer *BufferManager) UpdatePlayer(entityId string, value int16) error {
	isAddedToExistingPage, errTryToAppend := buffer.tryAppendToPage(entityId, value)
	if errTryToAppend != nil {
		return errTryToAppend
	}

	if !isAddedToExistingPage {
		pageAdress := buffer.addPage(entityId, value)

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

func (buffer *BufferManager) addPage(entityId string, value int16) *Page {
	buffer.pages = append(buffer.pages, newPage(entityId, value))
	return &buffer.pages[len(buffer.pages)-1]
}

func (buffer *BufferManager) tryAppendToPage(entityId string, value int16) (bool, error) {
	pageAdresses, err := buffer.getPageAdresses(entityId)
	if err != nil {
		return false, err
	}

	for i := 0; i < len(pageAdresses); i++ {
		if !pageAdresses[i].isPageFull() {
			pageAdresses[i].append(value)
			return true, nil
		}
	}
	return false, nil
}
