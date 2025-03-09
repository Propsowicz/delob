package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const clientKeySalt string = "Client Key"
const serverKeySalt string = "Server Key"

func addClientFirstAuthString(user string, c_nonce int) string {
	var auth string
	auth = addToAuth(auth, fmt.Sprintf("user=%s,", user))
	auth = addToAuth(auth, fmt.Sprintf("c_nonce=%d,", c_nonce))
	return auth
}

func addToAuth(auth string, s interface{}) string {
	var toAdd string
	switch v := s.(type) {
	case string:
		toAdd = v
	case int, int8:
		toAdd = fmt.Sprintf("%d", v)
	}

	return auth + toAdd
}

func parseClientFirstMessage(s string) (string, int, error) {
	parts := strings.Split(s, ",")
	const userPrefix string = "user="
	const cnoncePrefix string = "c_nonce="
	if parts[0][0:len(userPrefix)] == userPrefix && parts[1][0:len(cnoncePrefix)] == cnoncePrefix {
		if nonce, err := strconv.Atoi(parts[1][len(cnoncePrefix):]); err == nil {
			return parts[0][len(userPrefix):], nonce, nil
		}
	}

	return "", 0, fmt.Errorf("cannot parse client first message")
}

func generateNonce() int {
	return rand.Intn(256)
}

func computeHmacHash(arg_1, arg_2 []byte) []byte {
	mac := hmac.New(sha256.New, arg_1)
	mac.Write([]byte(arg_2))
	return mac.Sum(nil)
}

func computeSha256Hash(arg_1 []byte) []byte {
	hash := sha256.Sum256(arg_1)
	return hash[:]
}

func xorBytes(k, j []byte) []byte {
	if len(k) != len(j) {
		panic("byte slices must be of equal length")
	}
	result := make([]byte, len(k))
	for i := range k {
		result[i] = k[i] ^ j[i]
	}
	return result
}

func calculateHashedPassword(password string, salt []byte, iterations int) []byte {
	const keyLength int = 32

	hashedPassword := pbkdf2.Key([]byte(password), salt, iterations, keyLength, sha256.New)
	return hashedPassword
}

func generateRandomHash() []byte {
	randomBytes := make([]byte, 32)
	hash := sha256.Sum256(randomBytes)

	return hash[:]
}
