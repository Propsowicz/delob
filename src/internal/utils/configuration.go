package utils

import (
	"os"
)

const PAGE_SIZE int16 = 512
const DOCKER_ENV string = "docker"

func DockerEnvironment() bool {
	buildEnv := os.Getenv("BUILD_ENV")
	return buildEnv == DOCKER_ENV
}
