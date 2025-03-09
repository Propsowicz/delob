package auth

import (
	"fmt"
	"testing"
)

func Test_IfCanParseClientFirstMessageCorrectly(t *testing.T) {
	authM := NewAuthenticationManager()
	userMock := "testUser"

	user, _, _, err := authM.ParseClientFirstMessageToAuthString(fmt.Sprintf("user=%s,c_nonce=23", userMock))

	if err != nil {
		t.Errorf("Expected no error.")
	}
	if user != userMock {
		t.Errorf("Cannot parse user.")
	}
}
