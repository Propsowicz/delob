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

// func GenerateHashedPassword(password, salt string, iterations int) string {
// 	const keyLength int = 32

// 	hashedPassword := pbkdf2.Key([]byte(password), salt, iterations, keyLength, sha256.New)
// 	return hashedPassword
// }
