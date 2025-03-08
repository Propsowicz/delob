package auth

import (
	"delob/internal/utils/logger"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"time"
)

type userData struct {
	User       string
	Salt       []byte
	Iterations int
	Hashed_pwd []byte
	Stored_key []byte
	Client_key []byte
}

type userDataCollection struct {
	userData []userData
}

func AddUser(user, password string) error {
	userData, err := ReadMetaData()
	if err != nil {
		return err
	}

	*userData = append(*userData, createNewUser(user, password))

	errWrite := WriteMetaData(userData)
	if errWrite != nil {
		return errWrite
	}
	return nil
}

func LoadUserData(user string) (userData, error) {
	users, err := ReadMetaData()
	if err != nil {
		return userData{}, err
	}

	for i := range *users {
		if (*users)[i].User == user {
			return (*users)[i], nil
		}
	}
	return userData{}, fmt.Errorf("cannot find a user with given name: %s", user)
}

func createNewUser(user, password string) userData {
	u := userData{
		User: user,
	}
	u.Salt = generateRandomHash()
	u.Iterations = generateNonce()
	u.Hashed_pwd = calculateHashedPassword(password, u.Salt, u.Iterations)
	u.Client_key = computeHmacHash(u.Hashed_pwd, []byte(clientKeySalt))
	u.Stored_key = computeSha256Hash(u.Client_key)

	return u
}

func catalog() string {
	var path string = ".auth"

	err := os.MkdirAll(path, 0644)
	if err != nil {
		logger.Error("", err)
	}
	return path
}

func path() string {
	return fmt.Sprintf("%s/users.delob", catalog())
}

func ReadMetaData() (*[]userData, error) {
	var path string = path()
	var result []userData

	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return &result, nil
	}
	defer f.Close()

	if err != nil {
		logger.Error("", err)
		return nil, err
	}

	decoder := gob.NewDecoder(f)
	if err := decoder.Decode(&result); err != nil {
		logger.Error("", err)
		return &result, err
	}
	return &result, nil
}

func WriteMetaData(headers *[]userData) error {
	var path string = path()
	tempPath := fmt.Sprintf("%s_temp.%d", path, time.Now().Unix())

	f, err := os.OpenFile(tempPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Error("", err)
		return err
	}

	defer func() {
		f.Close()
		if err != nil {
			os.Remove(tempPath)
		}
	}()

	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(headers); err != nil {
		logger.Error("", err)
	}

	if err = f.Sync(); err != nil {
		logger.Error("", err)
		return err
	}

	f.Close()
	return os.Rename(tempPath, path)
}
