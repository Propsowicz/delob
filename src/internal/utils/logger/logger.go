package logger

import (
	"fmt"
	"log"
)

const (
	InfoLevel  string = "INFO"
	ErrorLevel        = "ERROR"
)

func Error(traceId string, err error) {
	log.Print(universalLogFormat(traceId, err.Error(), ErrorLevel))
}

func Info(traceId, msg string) {
	log.Print(universalLogFormat(traceId, msg, InfoLevel))
}

func universalLogFormat(traceId, msg, level string) string {
	return fmt.Sprintf(" - %s [%s] %s",
		traceId,
		level,
		msg)
}
