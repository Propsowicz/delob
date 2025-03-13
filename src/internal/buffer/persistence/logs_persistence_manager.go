package persistence

import (
	"bytes"
	"delob/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"sync"
)

var logsSeparator []byte = []byte("\n")

type LogsPersistenceManager struct {
	syncMutex      sync.Mutex
	LogsFileExists bool
	fileName       string
	dir            string
	path           string
}

func NewLogsPersistenceManager() (*LogsPersistenceManager, error) {
	b := LogsPersistenceManager{
		fileName: "logs",
		dir:      ".data",
	}
	b.path = fmt.Sprintf("%s/%s.delob", b.dir, b.fileName)

	err := os.MkdirAll(b.dir, 0755)
	if err != nil {
		return &b, err
	}

	if _, err := os.Stat(b.path); !errors.Is(err, os.ErrNotExist) {
		b.LogsFileExists = true
	}

	return &b, nil
}

type Log struct {
	Ver      string
	AddedOn  int64
	ExprType string
	Expr     string
}

func NewDataLog(traceId, parsedExpressionType, parsedExpression string) Log {
	return Log{
		Ver:      "00",
		AddedOn:  utils.Timestamp(),
		ExprType: parsedExpressionType,
		Expr:     parsedExpression,
	}
}

func (b *LogsPersistenceManager) Read() ([]Log, error) {
	f, err := os.ReadFile(b.path)
	if err != nil {
		return nil, err
	}

	jsonLogs := bytes.Split(f, logsSeparator)
	result := []Log{}

	for i := 0; i < len(jsonLogs)-1; i++ {
		obj := Log{}
		err := json.Unmarshal(jsonLogs[i], &obj)
		if err != nil {
			return nil, err
		}

		result = append(result, obj)
	}
	return result, nil
}

func (b *LogsPersistenceManager) Append(log Log) error {
	activePath, logFileExists := b.getActivePath()

	byteLog, bufferLogChecksum, err := b.getLogData(log)

	b.syncMutex.Lock()
	if logFileExists {
		if err := createBackupCopy(b.path, activePath); err != nil {
			return err
		}
	}

	err = b.appendToFile(byteLog, activePath)
	if err != nil {
		return err
	}

	logAppendedSuccesfully := b.isLogSuccessfullyAppended(bufferLogChecksum, activePath)
	logFileIntegrityError := b.handleLogFileIntegrity(logAppendedSuccesfully, logFileExists, activePath)
	b.syncMutex.Unlock()

	return logFileIntegrityError
}

func (b *LogsPersistenceManager) appendToFile(byteLog []byte, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
	}()

	if _, err := f.Write([]byte(append(byteLog, logsSeparator...))); err != nil {
		return err
	}

	if err = f.Sync(); err != nil {
		return err
	}

	return f.Close()
}

func (b *LogsPersistenceManager) handleLogFileIntegrity(logAppendedSuccesfully, logFileExists bool, tempPath string) error {
	if !logFileExists {
		if logAppendedSuccesfully {
			return nil
		} else {
			return os.Remove(b.path)
		}
	} else {
		if logAppendedSuccesfully {
			return os.Rename(tempPath, b.path)
		} else {
			return os.Remove(tempPath)
		}
	}
}

func (b *LogsPersistenceManager) isLogSuccessfullyAppended(bufferLogChecksum uint32, path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return false
	}
	size := stat.Size()

	var offset int64 = 0
	if size > 1024 {
		offset = size - 1024
	}

	_, err = f.Seek(offset, io.SeekStart)
	if err != nil {
		return false
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	if err != nil {
		return false
	}

	lastLogs := bytes.Split(buf.Bytes(), logsSeparator)
	lastLog := lastLogs[len(lastLogs)-2]

	fileLogChecksum, err := utils.Calculate(string(lastLog))
	if err != nil {
		return false
	}
	return bufferLogChecksum == fileLogChecksum
}

func (b *LogsPersistenceManager) getActivePath() (string, bool) {
	logFileExists := exists(b.path)
	if !logFileExists {
		return b.path, false
	}
	rnd := rand.Intn(int(math.MaxInt32))
	tempPath := fmt.Sprintf("%s_back.%d", b.path, rnd)

	if exists(tempPath) {
		return b.getActivePath()
	}
	return tempPath, true
}

func (b *LogsPersistenceManager) getLogData(log Log) ([]byte, uint32, error) {
	jsonLog, err := json.Marshal(log)
	if err != nil {
		return []byte{}, 0, err
	}

	bufferLogChecksum, err := utils.Calculate(string(jsonLog))
	if err != nil {
		return []byte{}, 0, err
	}
	return jsonLog, bufferLogChecksum, nil
}

func createBackupCopy(original, temp string) error {
	sourceFile, err := os.Open(original)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(temp)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
