package main

import (
	"fmt"
	"os"
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
	fmt.Println("TODO: hashing", password)
}

func compare(password, hash string) {
	fmt.Println("TODO: comparing", password, hash)
}
