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

var testCard = Card{
	Question: "Что есть оториноларинголог?",
	Options:  []string{"Ухо-горло-нос", "Печень-желчь-кишка", "Глаза-язык-легкие"},
}

var templates = template.Must(template.ParseFiles("card.html"))

func main() {
	http.HandleFunc("/", makeHandler(viewHandler))
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
		fn(w, r, testCard)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, c *Card) {
	err := templates.ExecuteTemplate(w, tmpl+".html", c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
