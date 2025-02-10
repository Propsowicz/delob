package buffer

import (
	"fmt"
	"os"
	"strings"
)

// type BackupManager interface {
// 	Save()
// 	Load()
// }

type BackupManager struct {
	dataCatalogFileName string
	dataDirectory       string
	path                string
}

func NewBackupManager() (BackupManager, error) {
	b := BackupManager{
		dataCatalogFileName: "data",
		dataDirectory:       "backup",
	}
	b.path = fmt.Sprintf("%s/%s.db", b.dataDirectory, b.dataCatalogFileName)

	err := os.MkdirAll(b.dataDirectory, 0755)
	if err != nil {
		return b, err
	}

	return b, nil
}

type SimpleDataFormat struct {
	Pages   []Page
	Matches []Match
}

func (b *BackupManager) Read() ([]string, error) {
	f, err := os.ReadFile(b.path)
	if err != nil {
		return nil, err
	}

	expressions := strings.Split(string(f), ";")
	return expressions, nil
}

func (b *BackupManager) Append(elo string) error {
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
