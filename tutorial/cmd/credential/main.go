package main

import (
	"fmt"
	"graphql/tutorial/internal/jwt"
	"os"
)

func main() {

	fmt.Println(os.Args)
	if len(os.Args) < 3 {
		fmt.Println("Please provide user and password")
		return
	}

	user := os.Args[1]
	password := os.Args[2]
	fmt.Printf("User :%s, Password: %s", user, password)

	salt := jwt.GenerateRandomSalt(8)
	userhash := user + password + salt

	fmt.Printf("\n Salt:%s, Hash:%s", salt, jwt.HashString(userhash))
}
