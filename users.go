package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type Credentials struct {
	Email    string `json:"EMAIL"`
	Password string `json:"PASSWORD"`
}

type Sex int

const (
	Unknown Sex = iota
	Male
	Female
)

var sexID = map[Sex]int{
	Male:   1,
	Female: 2,
}

type User struct {
	Id         int         `json:"PERSINFO_ID,omitempty"`
	Name       string      `json:"FIRSTNAME,omitempty"`
	Middlename string      `json:"MIDDLENAME,omitempty"`
	Surname    string      `json:"LASTNAME,omitempty"`
	Sex        Sex         `json:"SEX_ID,omitempty"`
	Auth       Credentials `json:"CREDENTIALS,omitempty"`
}

func (u User) getFullName() string {
	return fmt.Sprintf("%s %s %s", u.Surname, u.Name, u.Middlename)
}

func newUser(name, mname, surname, sex, email, password string) (u User) {
	s, _ := strconv.Atoi(sex)
	return User{
		Name:       name,
		Middlename: mname,
		Surname:    surname,
		Sex:        NewSex(s),
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

func NewSex(i int) Sex {
	for s, v := range sexID {
		if i == v {
			return s
		}
	}
	return Unknown
}
