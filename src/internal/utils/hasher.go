package utils

import (
	"crypto/sha256"
	"hash/fnv"

	"golang.org/x/crypto/pbkdf2"
)

func Calculate(id string) (uint32, error) {
	h := fnv.New32a()
	if _, err := h.Write([]byte(id)); err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}
func GenerateHashedPassword(password, salt string, iterations int) string {
	const keyLength int = 32

	hashedPassword := pbkdf2.Key([]byte(password), []byte(salt), iterations, keyLength, sha256.New)
	return string(hashedPassword)
}
