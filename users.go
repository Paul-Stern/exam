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
	id         int
	name       string
	middlename string
	surname    string
	auth       Credentials
}

type Users []User

var users = Users{
	User{
		id:         1,
		name:       "Евгений",
		middlename: "Семенович",
		surname:    "Коновалов",
		auth: Credentials{
			Email:    ***REMOVED***,
			Password: ***REMOVED***,
		},
	},
	User{
		id:         2,
		name:       "Юлиан",
		middlename: "Петрович",
		surname:    "Костоправ",
		auth: Credentials{
			Email:    ***REMOVED***,
			Password: ***REMOVED***,
		},
	},
	User{
		id:         3,
		name:       "Герман",
		middlename: "Станиславович",
		surname:    "Кривонос",
		auth: Credentials{
			Email:    ***REMOVED***,
			Password: ***REMOVED***,
		},
	},
}

func getUserById(id int) User {
	return users[id-1]
}

func getUserCreds(id int) Credentials {
	return users[id-1].auth
}

func getFullName(u User) string {
	return fmt.Sprintf("%s %s %s", u.surname, u.name, u.middlename)
}

func getUserByEmail(e string) (user User, err error) {
	for _, u := range users {
		if e == u.auth.Email {
			return u, nil
		}
	}
	err = errors.New("User not found")
	return User{}, err
}
