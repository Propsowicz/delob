package buffer

import (
	hasher "delob/internal/utils"
	"errors"
)

type PageDictionary struct {
	pagesData []PageData
}

type PageData struct {
	entityId     uint32
	pageAdresses []*Page
}

func (buffer *BufferManager) addPageToDictionary(entityId string, pageAdress *Page) error {
	hashedEntityId, err := hasher.Calculate(entityId)
	if err != nil {
		return err
	}

	newPageData := PageData{
		entityId:     hashedEntityId,
		pageAdresses: []*Page{pageAdress},
	}
	buffer.pageDictionary.pagesData = append(buffer.pageDictionary.pagesData, newPageData)
	return nil
}

func (buffer *BufferManager) appendPageToExistingKey(entityId string, pageAdress *Page) error {
	hashedEntityId, err := hasher.Calculate(entityId)
	if err != nil {
		return err
	}

	for i := range buffer.pageDictionary.pagesData {
		if buffer.pageDictionary.pagesData[i].entityId == hashedEntityId {
			buffer.pageDictionary.pagesData[i].pageAdresses =
				append(buffer.pageDictionary.pagesData[i].pageAdresses, pageAdress)
		}
	}

	return nil
}

func (buffer *BufferManager) getPageAdresses(entityId string) ([]*Page, error) {
	hashedEntityId, err := hasher.Calculate(entityId)
	if err != nil {
		return nil, err
	}

	for i := range buffer.pageDictionary.pagesData {
		if buffer.pageDictionary.pagesData[i].entityId == hashedEntityId {
			return buffer.pageDictionary.pagesData[i].pageAdresses, nil
		}
	}
	return nil, errors.New("cannot find entity with given id")
}
