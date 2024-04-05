package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	switch os.Args[1] {
	case "hash":
		if len(os.Args) != 3 {
			printUsage()
			return
		}
		hash(os.Args[2])
	case "compare":
		if len(os.Args) != 4 {
			printUsage()
			return
		}
		compare(os.Args[2], os.Args[3])
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println(os.Args[0], "usage:")
	fmt.Println("\thash: Hashes the password provided - hash <your password>")
	fmt.Println("\tcompare: Compares a password to a hash - compare <your password> <your hash>")
}

func hash(password string) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("error hashing: %v\n", password)
		return
	}
	fmt.Println(string(hashedBytes))
}

func compare(password, hash string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Invalid password:", password)
		return
	}
	fmt.Println("Password is correct!")
}
