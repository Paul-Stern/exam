package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Card struct {
	Question string
	Options  []string
}

type Credentials struct {
	email    string
	password string
}

type session struct {
	email  string
	expiry time.Time
}

var sessions = map[string]session{}

type User struct {
	id         int
	name       string
	middlename string
	surname    string
	auth       Credentials
}

type Users []User

type Test []Card

var cardOne = Card{
	Question: "Что есть оториноларинголог?",
	Options:  []string{"Ухо-горло-нос", "Печень-желчь-кишка", "Глаза-язык-легкие"},
}

var users = Users{
	User{
		id:         1,
		name:       "Евгений",
		middlename: "Семенович",
		surname:    "Коновалов",
		auth: Credentials{
			email:    ***REMOVED***,
			password: ***REMOVED***,
		},
	},
	User{
		id:         2,
		name:       "Юлиан",
		middlename: "Петрович",
		surname:    "Костоправ",
		auth: Credentials{
			email:    ***REMOVED***,
			password: ***REMOVED***,
		},
	},
	User{
		id:         3,
		name:       "Герман",
		middlename: "Станиславович",
		surname:    "Кривонос",
		auth: Credentials{
			email:    ***REMOVED***,
			password: ***REMOVED***,
		},
	},
}

var testOne = Test{
	cardOne,
	Card{
		Question: "Какой глаз ведущий у правши?",
		Options:  []string{"Левый", "Правый", "Средний (третий)"},
	},
}

var templates = template.Must(template.ParseFiles("test.html"))

func main() {
	http.HandleFunc("/test", makeHandler(viewHandler))
	http.HandleFunc("/login", signInHandler)
	log.Println("Server started. Listening to localhost:***REMOVED***")
	log.Fatal(http.ListenAndServe(":***REMOVED***", nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request, t Test) {
	switch r.Method {
	case "GET":
		renderTemplate(w, "test", &t)
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Printf("ParseForm() err: %v", err)
			return
		}
		f := r.PostForm
		log.Println(f)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "login.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Printf("ParseForm() err: %v", err)
			return
		}
		f := r.PostForm
		c := readCreds(f)
		u, err := getUserByEmail(c.email)
		if err != nil {
			log.Printf("getUserByEmail() err: %v", err)
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		if c == u.auth {
			log.Println("Login: succes. Redirecting...")
			http.Redirect(w, r, "/test", http.StatusFound)
		} else {
			err = errors.New("wrong password")
			log.Printf("Login: %v\n", err)
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}

func readCreds(f url.Values) Credentials {
	return Credentials{
		email:    f["email"][0],
		password: f["password"][0],
	}
}

func getUserCreds(id int) Credentials {
	return users[id-1].auth
}

func getUserById(id int) User {
	return users[id-1]
}

func getUserByEmail(e string) (user User, err error) {
	for _, u := range users {
		if e == u.auth.email {
			return u, nil
		}
	}
	err = errors.New("User not found")
	return User{}, err
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, Test)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, testOne)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, t *Test) {
	err := templates.ExecuteTemplate(w, tmpl+".html", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
