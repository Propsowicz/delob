package main

import (
	"delob/internal/auth"
	"flag"
	"fmt"
	"os"
)

func main() {
	initAppFlag := flag.Bool("init-app", false, "Initialiaztion app script. Provide username and password as positional arguments.")
	addUserFlag := flag.Bool("add-user", false, "Add a new user. Provide username and password as positional arguments.")
	loadUserFlag := flag.Bool("load-user", false, "Load user. Provide username as an argument.")

	flag.Parse()

	if *addUserFlag {
		shouldReturn := addUser()
		if shouldReturn {
			return
		}
	} else if *initAppFlag {
		shouldReturn := initApp()
		if shouldReturn {
			return
		}
	} else if *loadUserFlag {
		shouldReturn := loadUser()
		if shouldReturn {
			return
		}
	} else {
		fmt.Println("No action specified. Use --add-user flag to add a user.")
	}
}

func initApp() bool {
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: --init-app <username> <password>")
		os.Exit(1)
	}
	username := args[0]

	if userExists := userExists(username); userExists {
		return true
	}

	return addUser()
}

func addUser() bool {
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: --add-user <username> <password>")
		os.Exit(1)
	}
	username := args[0]
	password := args[1]

	if userExists := userExists(username); userExists {
		fmt.Printf("user %s already exists", username)
		return true
	}

	err := auth.AddUser(username, password)
	if err != nil {
		fmt.Println(err)
		return true
	}

	fmt.Printf("Adding user: %s\n", username)
	return false
}

func loadUser() bool {
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: --load-user <username>")
		os.Exit(1)
	}
	username := args[0]

	user, err := auth.LoadUserData(username)
	if err != nil {
		fmt.Println(err)
		return true
	}

	fmt.Printf("Loaded user: %s\n", user.User)
	return false
}

func userExists(username string) bool {
	_, err := auth.LoadUserData(username)
	return err == nil
}
