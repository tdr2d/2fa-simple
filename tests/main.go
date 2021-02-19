package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func main() {
	pass := "1234"

	fmt.Printf(HashPassword(pass))
}
