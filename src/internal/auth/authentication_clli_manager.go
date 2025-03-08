package auth

import (
	"delob/internal/utils/logger"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
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
	fmt.Println(userData)

	*userData = append(*userData, createNewUser(user, password))

	fmt.Println(userData)
	errWrite := WriteMetaData(userData)
	if errWrite != nil {
		return errWrite
	}

	return nil
}

func loadUserData(user string) userData {
	return createNewUser("myUsername", "myPassword")
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
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	decoder := gob.NewDecoder(f)
	if err := decoder.Decode(&result); err != nil {
		fmt.Println("Error decoding data:", err)
		return &result, err
	}

	return &result, nil
}

func WriteMetaData(headers *[]userData) error {
	var path string = path()
	tempPath := fmt.Sprintf("%s_temp.%d", path, time.Now().Unix())

	f, err := os.OpenFile(tempPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer func() {
		fmt.Println("????")
		f.Close()
		fmt.Println("!!")
		fmt.Println(err)
		if err != nil {
			errr := os.Remove(tempPath)
			fmt.Println(errr)
		}
	}()

	// jsonData, err := json.Marshal(headers)
	// if err != nil {
	// 	return err
	// }
	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(headers); err != nil {
		fmt.Println("Error encoding data:", err)
	}

	if err = f.Sync(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
	return os.Rename(tempPath, path)
}
