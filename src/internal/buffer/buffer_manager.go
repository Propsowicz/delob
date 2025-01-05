package buffer

import (
	"errors"
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
	if _, err := buffer.getPageAdresses(entityId); err == nil {
		return errors.New("player already exists")
	}

	buffer.syncMutex.Lock()

	pageAdress := buffer.addPage(entityId)
	buffer.addPageToDictionary(entityId, pageAdress)

	buffer.syncMutex.Unlock()
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
