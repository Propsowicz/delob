package auth

import (
	"testing"
)

func Test_IfCannotCreateUsernameWithCommaSeparator(t *testing.T) {
	_, err := createNewUser("cannotCreateThisUser,", "password")

	if err == nil {
		t.Errorf("Expected error.")
	}
}

func Test_IfCannotCreateShortPassword(t *testing.T) {
	_, err := createNewUser("user", "abcde")

	if err == nil {
		t.Errorf("Expected error.")
	}
}
