package buffer

import (
	hasher "delob/internal/utils"
	"fmt"
)

type PageDictionary struct {
	pagesData []*PageData
}

type PageData struct {
	hashedEntityId    uint32
	entityId          string
	pageAdresses      []*Page
	transactionStatus transactionStatus
}

func (buffer *BufferManager) addPageToDictionary(entityId string, pageAdress *Page) (*PageData, error) {
	hashedEntityId, err := hasher.Calculate(entityId)
	if err != nil {
		return nil, err
	}

	newPageData := PageData{
		hashedEntityId:    hashedEntityId,
		entityId:          entityId,
		pageAdresses:      []*Page{pageAdress},
		transactionStatus: inProgress,
	}
	buffer.pageDictionary.pagesData = append(buffer.pageDictionary.pagesData, &newPageData)
	return buffer.pageDictionary.pagesData[len(buffer.pageDictionary.pagesData)-1], nil
}

func (buffer *BufferManager) appendPageToExistingId(entityId string, pageAdress *Page) error {
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

type transactionStatusCondition func(transactionStatus transactionStatus) bool

func isInProgressOrSuccess(transactionStatus transactionStatus) bool {
	return transactionStatus == success || transactionStatus == inProgress
}

func isSuccess(transactionStatus transactionStatus) bool {
	return transactionStatus == success
}

func (buffer *BufferManager) getPageAdresses(entityId string,
	transactionStatusCondition transactionStatusCondition) ([]*Page, error) {
	hashedEntityId, err := hasher.Calculate(entityId)
	if err != nil {
		return nil, err
	}

	for i := range buffer.pageDictionary.pagesData {
		if buffer.pageDictionary.pagesData[i].hashedEntityId == hashedEntityId {
			if buffer.pageDictionary.pagesData[i].transactionStatus == failed {
				continue
			}

			if transactionStatusCondition(buffer.pageDictionary.pagesData[i].transactionStatus) {
				return buffer.pageDictionary.pagesData[i].pageAdresses, nil
			}
			return nil, fmt.Errorf("cannot find entity with given id: %s", entityId)
		}
	}
	return nil, fmt.Errorf("cannot find entity with given id: %s", entityId)
}
