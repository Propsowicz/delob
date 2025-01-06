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

func (buffer *BufferManager) AddPlayer(entityId string) error {
	if _, err := buffer.getPageAdresses(entityId); err == nil {
		return fmt.Errorf("player already exists: %s", entityId)
	}

	buffer.syncMutex.Lock()

	pageAdress := buffer.addPage(entityId)
	err := buffer.addPageToDictionary(entityId, pageAdress)
	if err != nil {
		return err
	}

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
