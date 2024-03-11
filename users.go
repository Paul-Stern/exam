package main

import (
	"fmt"
	"net/http"
)

type Credentials struct {
	Email    string `json:"EMAIL"`
	Password string `json:"PASSWORD"`
}

type User struct {
	Id         int         `json:"PERSINFO_ID,omitempty"`
	Name       string      `json:"FIRSTNAME,omitempty"`
	Middlename string      `json:"MIDDLENAME,omitempty"`
	Surname    string      `json:"LASTNAME,omitempty"`
	Sex        string      `json:"SEX_ID,omitempty"`
	Auth       Credentials `json:"CREDENTIALS,omitempty"`
}

func (u User) getFullName() string {
	return fmt.Sprintf("%s %s %s", u.Surname, u.Name, u.Middlename)
}

func newUser(name, mname, surname, email, password string) (u User) {
	return User{
		Name:       name,
		Middlename: mname,
		Surname:    surname,
		Auth: Credentials{
			Email:    email,
			Password: password,
		},
	}
}

// Get user from REST server
func getUser(email string) (u User, err error) {
	u.Auth.Email = email
	query := "?email=" + email
	url := baseUrl(cfg) + "/persons" + query
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	message, err := read[User](resp)
	if err != nil {
		return
	}
	u = message.Data
	return

}
