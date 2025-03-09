package main

import (
	"delob/internal/auth"
	"flag"
	"fmt"
	"os"
)

func main() {
	addUser := flag.Bool("add-user", false, "Add a new user. Provide username and password as positional arguments.")
	loadUser := flag.Bool("load-user", false, "Load user. Provide username as an argument.")

	flag.Parse()

	if *addUser {
		args := flag.Args()
		if len(args) < 2 {
			fmt.Println("Usage: --add-user <username> <password>")
			os.Exit(1)
		}
		username := args[0]
		password := args[1]

		err := auth.AddUser(username, password)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Adding user: %s\n", username)
	} else if *loadUser {
		args := flag.Args()
		if len(args) != 1 {
			fmt.Println("Usage: --load-user <username>")
			os.Exit(1)
		}
		username := args[0]

		user, err := auth.LoadUserData(username)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Loaded user: %s\n", user.User)
	} else {
		fmt.Println("No action specified. Use --add-user flag to add a user.")
	}
}
