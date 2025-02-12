package buffer

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// type BackupManager interface {
// 	Save()
// 	Load()
// }

// TODO rename it
type LogDataManager struct {
	IsBackupFileExists  bool
	dataCatalogFileName string
	dataDirectory       string
	path                string
}

func NewLogDataManager() (LogDataManager, error) {
	b := LogDataManager{
		dataCatalogFileName: "dict",
		dataDirectory:       "log_data",
	}
	b.path = fmt.Sprintf("%s/%s.delob", b.dataDirectory, b.dataCatalogFileName)

	err := os.MkdirAll(b.dataDirectory, 0755)
	if err != nil {
		return b, err
	}

	if _, err := os.Stat(b.path); !errors.Is(err, os.ErrNotExist) {
		b.IsBackupFileExists = true
	}

	return b, nil
}

type SimpleDataFormat struct {
	Pages   []Page
	Matches []Match
}

func (b *LogDataManager) Read() ([]string, error) {
	f, err := os.ReadFile(b.path)
	if err != nil {
		return nil, err
	}
	expressions := strings.Split(string(f), ";")
	result := []string{}

	for i := 0; i < len(expressions)-1; i++ {
		result = append(result, expressions[i]+";")
	}
	return result, nil
}

func (b *LogDataManager) Append(elo string) error {
	f, err := os.OpenFile(b.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer func() {
		f.Close()
	}()

	if _, err := f.Write([]byte(elo)); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	return nil
}
