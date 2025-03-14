package utils

import "time"

func Timestamp() int64 {
	return time.Now().UnixMilli()
}

func TimestampMinutesOffset(minutes int8) int64 {
	return time.Now().Add(time.Duration(minutes) * time.Minute).UnixMilli()
}
