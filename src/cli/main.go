package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define the flag. When --add-user is provided, it will be true.
	addUser := flag.Bool("add-user", false, "Add a new user. Provide username and password as positional arguments.")

	// Parse command-line flags
	flag.Parse()

	// If the --add-user flag is set, expect two additional arguments: username and password.
	if *addUser {
		args := flag.Args() // returns the positional arguments that are not flags
		if len(args) < 2 {
			fmt.Println("Usage: --add-user <username> <password>")
			os.Exit(1)
		}
		username := args[0]
		password := args[1]
		fmt.Printf("Adding user: %s with password: %s\n", username, password)
	} else {
		fmt.Println("No action specified. Use --add-user flag to add a user.")
	}
}
