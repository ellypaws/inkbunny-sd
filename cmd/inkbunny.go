package main

import (
	"fmt"
	"github.com/ellypaws/inkbunny/api"
	"golang.org/x/term"
	"log"
	"os"
)

func login(u *api.Credentials) error {
	if u == nil {
		return fmt.Errorf("nil user")
	}
	if u.Sid == "" {
		user, err := loginPrompt().Login()
		if err != nil {
			return err
		}
		*u = *user
	}
	return nil
}

func loginPrompt() *api.Credentials {
	var user api.Credentials
	fmt.Print("Enter username [guest]: ")
	fmt.Scanln(&user.Username)
	fmt.Print("Enter password: ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(bytePassword)

	return &user
}

func logout(u *api.Credentials) {
	err := u.Logout()
	if err != nil {
		log.Fatalf("error logging out: %v", err)
	}
}
