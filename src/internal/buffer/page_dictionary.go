package buffer

import (
	hasher "delob/internal/utils"
)

type PageDictionary struct {
	pagesData []*PageData
}

type PageData struct {
	hashedKey         uint32
	key               string
	pageAdresses      []*Page
	transactionStatus transactionStatus
}

func (buffer *BufferManager) addPageToDictionary(key string, pageAdress *Page) (*PageData, error) {
	hashedKey, err := hasher.Calculate(key)
	if err != nil {
		return nil, err
	}

	newPageData := PageData{
		hashedKey:         hashedKey,
		key:               key,
		pageAdresses:      []*Page{pageAdress},
		transactionStatus: inProgress,
	}
	buffer.pageDictionary.pagesData = append(buffer.pageDictionary.pagesData, &newPageData)
	return buffer.pageDictionary.pagesData[len(buffer.pageDictionary.pagesData)-1], nil
}

func (buffer *BufferManager) appendPageToExistingId(key string, pageAdress *Page) error {
	hashedKey, err := hasher.Calculate(key)
	if err != nil {
		return err
	}

	for i := range buffer.pageDictionary.pagesData {
		if buffer.pageDictionary.pagesData[i].hashedKey == hashedKey {
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

func (buffer *BufferManager) getPageAdresses(key string,
	transactionStatusCondition transactionStatusCondition) ([]*Page, error) {
	hashedKey, err := hasher.Calculate(key)
	if err != nil {
		return nil, err
	}

	for i := range buffer.pageDictionary.pagesData {
		if buffer.pageDictionary.pagesData[i].hashedKey == hashedKey {
			if buffer.pageDictionary.pagesData[i].transactionStatus == failed {
				continue
			}

			if transactionStatusCondition(buffer.pageDictionary.pagesData[i].transactionStatus) {
				return buffer.pageDictionary.pagesData[i].pageAdresses, nil
			}
			return nil, errorPlayerDoesNotExists(key)
		}
	}
	return nil, errorPlayerDoesNotExists(key)
}
