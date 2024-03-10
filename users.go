package main

import (
	"errors"
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
	Auth       Credentials `json:"CREDENTIALS,omitempty"`
}

type Users []User

var users = Users{
	User{
		Id:         1,
		Name:       "Евгений",
		Middlename: "Семенович",
		Surname:    "Коновалов",
		Auth: Credentials{
			Email:    ***REMOVED***,
			Password: ***REMOVED***,
		},
	},
	User{
		Id:         2,
		Name:       "Юлиан",
		Middlename: "Петрович",
		Surname:    "Костоправ",
		Auth: Credentials{
			Email:    ***REMOVED***,
			Password: ***REMOVED***,
		},
	},
	User{
		Id:         3,
		Name:       "Герман",
		Middlename: "Станиславович",
		Surname:    "Кривонос",
		Auth: Credentials{
			Email:    ***REMOVED***,
			Password: ***REMOVED***,
		},
	},
}

func getUserById(id int) User {
	return users[id-1]
}

func getUserCreds(id int) Credentials {
	return users[id-1].Auth
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

func getUserByEmail(e string) (user User, err error) {
	for _, u := range users {
		if e == u.Auth.Email {
			return u, nil
		}
	}
	err = errors.New("User not found")
	return User{}, err
}

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
