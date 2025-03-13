package buffer

import "fmt"

func errorPlayerAlreadyExists(key string) error {
	return fmt.Errorf("player with key '%s' already exists", key)
}

func errorPlayerDoesNotExists(key string) error {
	return fmt.Errorf("cannot find player with key '%s'", key)
}
