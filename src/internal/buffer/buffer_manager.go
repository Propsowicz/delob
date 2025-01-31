package buffer

import (
	"delob/internal/utils"
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

func (buffer *BufferManager) AddPlayer(entityId string) error {

	// move it to processor?
	if _, err := buffer.getPageAdresses(entityId); err == nil {
		return fmt.Errorf("player already exists: %s", entityId)
	}

	buffer.syncMutex.Lock()

	record := Record{
		Method: Add,
		// TODO move it out of here
		Value: utils.INITIAL_ELO,
	}

	pageAdress := buffer.addPage(entityId, record)
	err := buffer.addPageToDictionary(entityId, pageAdress)
	if err != nil {
		return err
	}

	buffer.syncMutex.Unlock()
	return nil
}

func (buffer *BufferManager) UpdatePlayer(entityId string, record Record) error {
	isAddedToExistingPage, errTryToAppend := buffer.tryAppendToPage(entityId, record)
	if errTryToAppend != nil {
		return errTryToAppend
	}

	if !isAddedToExistingPage {
		buffer.syncMutex.Lock()
		pageAdress := buffer.addPage(entityId, record)
		errAppendPageToExistingKey := buffer.appendPageToExistingKey(entityId, pageAdress)
		if errAppendPageToExistingKey != nil {
			return errAppendPageToExistingKey
		}

		buffer.syncMutex.Unlock()
	}

	return nil
}

func (buffer *BufferManager) GetPage(entityId string) ([]Page, error) {
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
