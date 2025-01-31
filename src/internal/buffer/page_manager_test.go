package buffer

// import (
// 	"testing"
// )

// func Test_IfTryAppendToPageButNoEntityCreated_ThenThrow(t *testing.T) {
// 	bufferMock := NewBufferManager()
// 	var entityIdMock string = "testId"

// 	result, err := bufferMock.tryAppendToPage(entityIdMock, Record{})

// 	if result == true {
// 		t.Errorf("Result should be false.")
// 	}
// 	if err == nil {
// 		t.Errorf("Result should return error.")
// 	}
// }

// func Test_IfTryAppendToPageThatIsNotfull_ThenSucced(t *testing.T) {
// 	bufferMock := NewBufferManager()
// 	var entityIdMock string = "testId"

// 	bufferMock.addPage(entityIdMock, Record{
// 		Method: Add,
// 		Value:  1500,
// 	})

// 	result, err := bufferMock.tryAppendToPage(entityIdMock, Record{
// 		Method: Add,
// 		Value:  25,
// 	})

// 	if result != true {
// 		t.Errorf("Result should be true.")
// 	}
// 	if err != nil {
// 		t.Errorf("Result should not return error.")
// 	}
// }
