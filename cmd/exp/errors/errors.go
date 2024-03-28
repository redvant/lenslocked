package main

import (
	"errors"
	"fmt"
)

func Connect() error {
	// try to connect
	// pretend we got an error
	// panic("connection failed")
	return errors.New("connection failed")
}

func CreateUser() error {
	err := Connect()
	if err != nil {
		// We can add more context here!
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func CreateOrg() error {
	err := CreateUser()
	if err != nil {
		return fmt.Errorf("create org: %w", err)
	}
	return nil
}

func main() {
	// err := CreateUser()
	// if err != nil {
	// log.Println(err)
	// }
	// err = CreateOrg()
	// if err != nil {
	// log.Println(err)
	// }

	err := B()
	// TODO: Determine if the err variable is an 'ErrNotFound'
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Error match")
	}
}

var ErrNotFound = errors.New("not found")

func A() error {
	return ErrNotFound
}

func B() error {
	err := A()
	if err != nil {
		return fmt.Errorf("b: %w", err)
	}
	return nil
}
