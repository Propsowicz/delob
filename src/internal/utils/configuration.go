package utils

import "time"

const PAGE_SIZE int16 = 512

func Timestamp() int64 {
	return time.Now().UnixMilli()
}

func TimestampMinutesOffset(minutes int8) int64 {
	return time.Now().Add(time.Duration(minutes) * time.Minute).UnixMilli()
}
