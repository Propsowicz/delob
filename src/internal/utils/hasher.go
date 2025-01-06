package utils

import (
	"hash/fnv"
)

func Calculate(id string) (uint32, error) {
	h := fnv.New32a()
	if _, err := h.Write([]byte(id)); err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}
