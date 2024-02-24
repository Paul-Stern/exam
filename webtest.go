package main

import (
	"html/template"
	"log"
	"net/http"
)

type Card struct {
	Question string
	Options  []string
}

type User struct {
	id         int
	name       string
	middlename string
	surname    string
	email      string
}

type Users []User

type Test []Card

var cardOne = Card{
	Question: "Что есть оториноларинголог?",
	Options:  []string{"Ухо-горло-нос", "Печень-желчь-кишка", "Глаза-язык-легкие"},
}

var userOne = User{
	id:         1,
	name:       "Евгений",
	middlename: "Семенович",
	surname:    "Коновалов",
	email:      "konoval@example.com",
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
	http.HandleFunc("/", makeHandler(viewHandler))
	log.Println("Server started. Listening to localhost:***REMOVED***")
	log.Fatal(http.ListenAndServe(":***REMOVED***", nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request, t Test) {
	renderTemplate(w, "test", &t)
	r.ParseForm()
	vals := r.Form
	log.Println(vals)
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
