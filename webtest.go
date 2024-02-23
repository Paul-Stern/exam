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

type test []Card

var testCard = Card{
	Question: "Что есть оториноларинголог?",
	Options:  []string{"Ухо-горло-нос", "Печень-желчь-кишка", "Глаза-язык-легкие"},
}

var testOne = test{
	Card{
		Question: "Что есть оториноларинголог?",
		Options:  []string{"Ухо-горло-нос", "Печень-желчь-кишка", "Глаза-язык-легкие"},
	},
	Card{
		Question: "Какой глаз ведущий у правши?",
		Options:  []string{"Левый", "Правый", "Срений (третий)"},
	},
}

var templates = template.Must(template.ParseFiles("card.html"))

func main() {
	http.HandleFunc("/", makeHandler(viewHandler))
	log.Println("Server started. Listening to localhost:***REMOVED***")
	log.Fatal(http.ListenAndServe(":***REMOVED***", nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request, c Card) {
	renderTemplate(w, "card", &c)
	r.ParseForm()
	vals := r.Form
	log.Println(vals)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, Card)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, c := range testOne {
			fn(w, r, c)
		}
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, c *Card) {
	err := templates.ExecuteTemplate(w, tmpl+".html", c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
