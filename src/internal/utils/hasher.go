package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

func GenerateKey() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func GenerateHashedPassword(password, salt string, iterations int) string {
	const keyLength int = 32

	hashedPassword := pbkdf2.Key([]byte(password), []byte(salt), iterations, keyLength, sha256.New)
	return string(hashedPassword)
}
