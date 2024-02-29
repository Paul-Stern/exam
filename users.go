package main

import (
	"errors"
	"fmt"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id         int         `json:"PERSINFO_ID"`
	Name       string      `json:"FIRSTNAME"`
	Middlename string      `json:"MIDDLENAME"`
	Surname    string      `json:"LASTNAME"`
	Auth       Credentials `json:"CREDENTIALS"`
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
