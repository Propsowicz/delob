package buffer

import (
	"errors"
	"fmt"
)

type PageDictionary struct {
	pagesData []PageData
}

type PageData struct {
	entityId     string
	pageAdresses []*Page
}

func (buffer *BufferManager) addPageToDictionary(entityId string, pageAdress *Page) {
	newPageData := PageData{
		entityId:     entityId,
		pageAdresses: []*Page{pageAdress},
	}
	buffer.pageDictionary.pagesData = append(buffer.pageDictionary.pagesData, newPageData)
}

func (buffer *BufferManager) getPageAdresses(entityId string) ([]*Page, error) {
	if len(buffer.pageDictionary.pagesData) == 0 {
		return nil, errors.New("there are no pages in page dictionary")
	}
	fmt.Println(entityId)
	for i := range buffer.pageDictionary.pagesData {
		if buffer.pageDictionary.pagesData[i].entityId == entityId {
			return buffer.pageDictionary.pagesData[i].pageAdresses, nil
		}
	}
	return nil, errors.New("cannot find entity with given id")
}
