package utils

import (
	"os"
	"time"
)

const PAGE_SIZE int16 = 512
const DOCKER_ENV string = "docker"

func DockerEnvironment() bool {
	buildEnv := os.Getenv("BUILD_ENV")
	return buildEnv == DOCKER_ENV
}

func Timestamp() int64 {
	return time.Now().UnixMilli()
}

func TimestampMinutesOffset(minutes int8) int64 {
	return time.Now().Add(time.Duration(minutes) * time.Minute).UnixMilli()
}
