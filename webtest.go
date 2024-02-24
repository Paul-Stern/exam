package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	Question string
	Options  []string
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

type Test struct {
	User  string
	Cards []Card
}

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

var testOne = Test{
	User: getFullName(getUserById(1)),
	Cards: []Card{
		cardOne,
		Card{
			Question: "Какой глаз ведущий у правши?",
			Options:  []string{"Левый", "Правый", "Средний (третий)"},
		},
	},
}

var templates = template.Must(template.ParseFiles("test.html"))

func main() {
	http.HandleFunc("/test", makeHandler(viewHandler))
	http.HandleFunc("/login", signInHandler)
	http.HandleFunc("/json", jsonHandler)
	log.Println("Server started. Listening to localhost:***REMOVED***")
	log.Fatal(http.ListenAndServe(":***REMOVED***", nil))
}

func newTest(u User) Test {
	return Test{
		User: getFullName(u),
		Cards: []Card{
			cardOne,
			Card{
				Question: "Какой глаз ведущий у правши?",
				Options:  []string{"Левый", "Правый", "Средний (третий)"},
			},
		},
	}
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(users[1].auth)
	if err != nil {
		log.Printf("Json error: %v", err)
	}
	w.Header().Add("Content-type", "application/json")
	fmt.Fprintf(w, "%s", data)
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
		u, err := getUserByEmail(c.Email)
		if err != nil {
			log.Printf("getUserByEmail() err: %v", err)
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		if c == u.auth {
			sessionToken := uuid.NewString()
			expiresAt := time.Now().Add(2 * time.Hour)

			sessions[sessionToken] = session{
				email:  c.Email,
				expiry: expiresAt,
			}
			log.Println("Login: success. Redirecting...")
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
		Email:    f["email"][0],
		Password: f["password"][0],
	}
}

func getUserCreds(id int) Credentials {
	return users[id-1].auth
}

func getUserById(id int) User {
	return users[id-1]
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

func makeHandler(fn func(http.ResponseWriter, *http.Request, Test)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t Test
		for _, s := range sessions {
			u, err := getUserByEmail(s.email)
			if err != nil {
				log.Printf("Get current user error: %v", err)
			}
			t = newTest(u)
		}
		fn(w, r, t)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, t *Test) {
	err := templates.ExecuteTemplate(w, tmpl+".html", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
