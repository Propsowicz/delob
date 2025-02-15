package buffer

import (
	"bytes"
	"delob/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// type BackupManager interface {
// 	Save()
// 	Load()
// }

var logsSeparator []byte = []byte{95}

// TODO rename it
type DataLogsDictionaryManager struct {
	IsLogsDictionaryFileExists bool
	dataCatalogFileName        string
	dataDirectory              string
	path                       string
}

func NewDataLogsDictionaryManager() (DataLogsDictionaryManager, error) {
	b := DataLogsDictionaryManager{
		dataCatalogFileName: "dict",
		dataDirectory:       "log_data",
	}
	b.path = fmt.Sprintf("%s/%s.delob", b.dataDirectory, b.dataCatalogFileName)

	err := os.MkdirAll(b.dataDirectory, 0755)
	if err != nil {
		return b, err
	}

	if _, err := os.Stat(b.path); !errors.Is(err, os.ErrNotExist) {
		b.IsLogsDictionaryFileExists = true
	}

	return b, nil
}

type DataLog struct {
	TraceId              string
	AddTimestamp         int64
	ParsedExpressionType string
	ParsedExpression     string
}

func NewDataLog(traceId, parsedExpressionType, parsedExpression string) DataLog {
	return DataLog{
		TraceId:              traceId,
		AddTimestamp:         utils.Timestamp(),
		ParsedExpressionType: parsedExpressionType,
		ParsedExpression:     parsedExpression,
	}
}

func (b *DataLogsDictionaryManager) Read() ([]DataLog, error) {
	f, err := os.ReadFile(b.path)
	if err != nil {
		return nil, err
	}

	jsonLogs := bytes.Split(f, logsSeparator)
	result := []DataLog{}

	for i := 0; i < len(jsonLogs)-1; i++ {
		obj := DataLog{}
		err := json.Unmarshal(jsonLogs[i], &obj)
		if err != nil {
			return nil, err
		}

		result = append(result, obj)
	}
	return result, nil
}

func (b *DataLogsDictionaryManager) Append(log DataLog) error {
	f, err := os.OpenFile(b.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer func() {
		f.Close()
	}()

	jsonLog, err := json.Marshal(log)
	if err != nil {
		return err
	}

	if _, err := f.Write([]byte(append(jsonLog, logsSeparator...))); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	return nil
}
