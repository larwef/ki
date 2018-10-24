package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
)

func main() {

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatalf("Error getting password input: %v", err)
	}
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error gettin password input: %v", err)
	}

	fmt.Printf("\nHashed password: %s\n", string(passwordHash))
}
