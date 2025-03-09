package auth

import (
	"encoding/hex"
	"fmt"
)

type AuthenticationManager struct {
	sessionManager sessionManager
}

func NewAuthenticationManager() AuthenticationManager {
	return AuthenticationManager{}
}

func (a *AuthenticationManager) IsUserAuthenticated(user, ip string) bool {
	return a.sessionManager.IsSessionValid(user, ip)
}

func (a *AuthenticationManager) AddServerFirstMessage(auth, user string) (string, error) {
	userData, err := LoadUserData(user)
	if err != nil {
		return "", err
	}

	auth = addToAuth(auth, fmt.Sprintf("s_nonce=%d,", generateNonce()))
	auth = addToAuth(auth, fmt.Sprintf("salt=%s,", hex.EncodeToString(userData.Salt)))
	auth = addToAuth(auth, fmt.Sprintf("iterations=%d", userData.Iterations))
	return auth, nil
}

func (a *AuthenticationManager) ParseClientFirstMessageToAuthString(message string) (string, int, string, error) {
	user, c_nonce, err := parseClientFirstMessage(message)
	if err != nil {
		return "", 0, "", err
	}

	return user, c_nonce, message, nil
}

func (a *AuthenticationManager) Verify(proof, user, ip, auth string) bool {
	userData, err := LoadUserData(user)
	if err != nil {
		return false
	}
	clientSignature := computeHmacHash(userData.Stored_key, []byte(auth))
	serverSideProof := hex.EncodeToString(xorBytes(userData.Client_key, clientSignature))
	isProofValid := serverSideProof == proof

	a.sessionManager.AddToSession(user, ip)
	return isProofValid
}
