package buffer

import (
	hasher "delob/internal/utils"
	"errors"
)

type PageDictionary struct {
	pagesData []PageData
}

type PageData struct {
	hashedEntityId uint32
	entityId       string
	pageAdresses   []*Page
}

func (buffer *BufferManager) addPageToDictionary(entityId string, pageAdress *Page) error {
	hashedEntityId, err := hasher.Calculate(entityId)
	if err != nil {
		return err
	}

	newPageData := PageData{
		hashedEntityId: hashedEntityId,
		entityId:       entityId,
		pageAdresses:   []*Page{pageAdress},
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
		if buffer.pageDictionary.pagesData[i].hashedEntityId == hashedEntityId {
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
		if buffer.pageDictionary.pagesData[i].hashedEntityId == hashedEntityId {
			return buffer.pageDictionary.pagesData[i].pageAdresses, nil
		}
	}
	return nil, errors.New("cannot find entity with given id")
}
