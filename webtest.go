package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Option struct {
	Id   int
	Text string
}

type Card struct {
	Id       int
	Question string
	Options  []Option
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type session struct {
	email  string
	expiry time.Time
}

type message []byte

//	type restBlock struct {
//		id             int    `json:"id"`
//		task_text      string `json:"task_text"`
//		task_answer_id int    `json:"task_answer_id"`
//		answer_text    string `json:"answer_text"`
//	}
type restBlock struct {
	Id        int       `json:"ID"`
	Task_text string    `json:"TASK_TEXT"`
	Answers   []restOpt `json:"ANSWERS"`
}

type restOpt struct {
	Id          int    `json:"ID"`
	Answer_text string `json:"ANSWER_TEXT"`
}

type restBlocks []restBlock

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

const (
	getUrl       = "http://***REMOVED***:***REMOVED******REMOVED******REMOVED***
	urlQuestGet  = "http://***REMOVED***:***REMOVED******REMOVED******REMOVED***"
	urlQuestPost = "http://***REMOVED***:***REMOVED******REMOVED******REMOVED***"
)

var cardOne = Card{
	Question: "Что есть оториноларинголог?",
	// Options:  []string{"Ухо-горло-нос", "Печень-желчь-кишка", "Глаза-язык-легкие"},
	Options: []Option{
		Option{
			Id:   9999,
			Text: "Ухо-горло-нос",
		},
		Option{
			Id:   10000,
			Text: "Печень-желчь-кишка",
		},
	},
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

// var testOne = Test{
// 	User: getFullName(getUserById(1)),
// 	Cards: []Card{
// 		cardOne,
// 		Card{
// 			Question: "Какой глаз ведущий у правши?",
// 			Options:  []string{"Левый", "Правый", "Средний (третий)"},
// 		},
// 	},
// }

var templates = template.Must(template.ParseFiles("test.html"))

func main() {
	http.HandleFunc("/test", makeHandler(viewHandler))
	http.HandleFunc("/login", signInHandler)
	http.HandleFunc("/json", jsonHandler)
	http.HandleFunc("/req", reqHandler)
	log.Println("Server started. Listening to localhost:***REMOVED***")
	log.Fatal(http.ListenAndServe(":***REMOVED***", nil))
}

func newTest(u User, c []Card) Test {
	// return Test{
	// 	User: getFullName(u),
	// 	Cards: []Card{
	// 		cardOne,
	// 		Card{
	// 			Question: "Какой глаз ведущий у правши?",
	// 			Options:  []string{"Левый", "Правый", "Средний (третий)"},
	// 		},
	// 	},
	// }
	return Test{
		User: getFullName(u),
		// Cards: []Card{
		// 	newCard(
		// 		9999,
		// 		"Что такое оториноларинголог?",
		// 		[]Option{
		// 			newOption(501, "Ухо-горло-нос"),
		// 			newOption(502, "Почки-глаза"),
		// 		},
		// 	),
		// },
		Cards: c,
	}
}

func newOption(id int, text string) Option {
	var o Option
	o.Id = id
	o.Text = text
	return o
}

func newCard(id int, question string, opts []Option) Card {
	var c Card
	c.Id = id
	c.Question = question
	c.Options = opts
	return c

}

func getData(url string) (m message, err error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("getUrl http.Get() error: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("getUrl io.ReadAll() error: %v", err)
	}
	if json.Valid(body) {
		return body, nil
	}
	err = errors.New("not a valid json")
	return message{}, err
}

/*
func getCards() []Card {
	var rbs restBlocks
	var cards []Card
	data, err := getData(urlQuestGet)
	if err != nil {
		log.Fatalf("get cards error: %v", err)
	}
	err = json.Unmarshal(data, &rbs)
	if err != nil {
		log.Fatalf("get cards error: %v", err)
	}
	var c Card
	for _, b := range rbs {
		if c.Id == b.id {
			c.Options = append(
				c.Options,
				newOption(b.task_answer_id, b.answer_text),
			)
		}
		if c.Id != b.id {
			if c.Id > 0 {
				cards = append(cards, c)
			}
			c.Id = b.id
			c.Question = b.task_text
		}
	}
	return cards
}
*/

func getCards() []Card {
	var rbs restBlocks
	var cards []Card
	data, err := getData(urlQuestGet)
	if err != nil {
		log.Fatalf("get cards error: %v", err)
	}
	err = json.Unmarshal(data, &rbs)
	if err != nil {
		log.Fatalf("get cards error: %v", err)
	}
	var c Card
	for _, block := range rbs {
		c.Id = block.Id
		c.Question = block.Task_text
		c.Options = []Option{}
		for _, o := range block.Answers {
			c.Options = append(
				c.Options,
				newOption(o.Id, o.Answer_text),
			)
		}
		cards = append(cards, c)
	}
	return cards
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getData(urlQuestGet)
	if err != nil {
		log.Fatalf("req error: %v\n", err)
	}
	w.Header().Add("Content-type", "application/json")
	fmt.Fprintf(w, "%s", data)
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
			log.Fatalf("ParseForm() err: %v", err)
		}
		// f := r.PostForm
		c := readCreds(r.PostForm)
		u, err := getUserByEmail(c.Email)
		if err != nil {
			log.Printf("getUserByEmail() err: %v", err)
			http.Redirect(w, r, "/login", http.StatusFound)
		}

		/*
			form, err := json.Marshal(c)
			if err != nil {
				log.Printf("Marshal error: %v", err)
			}

			resp, err := http.Post(
				urlQuestPost,
				"application/json",
				bytes.NewBuffer(form),
			)
			log.Printf("json sent")
			if err != nil {
				log.Fatalf("login error: %v", err)
			}
			// Read body
			got, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("login error: %v", err)
			}
			fmt.Fprintf(w, "%s", getCards())
			if json.Valid(got) {
				r.Header.Add("Content-Type", "application/json")
				fmt.Fprintf(w, "%s", got)
				break
			}
			fmt.Fprintf(w, "%s", got)
			log.Fatalln("got invalid json")
			log.Println("Login: success. Redirecting...")
			http.Redirect(w, r, "/test", http.StatusFound)
		*/

		// Check if credentials are correct
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

func getSession(u User) session {
	for uuid := range sessions {
		s := sessions[uuid]
		log.Printf("%v\n", s.email)
		return s
	}
	return session{}
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
		c := getCards()
		for _, s := range sessions {
			u, err := getUserByEmail(s.email)
			if err != nil {
				log.Printf("Get current user error: %v", err)
			}
			t = newTest(u, c)
		}
		// t = newTest(
		// 	getUserById(1),
		// 	c,
		// )
		fn(w, r, t)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, t *Test) {
	err := templates.ExecuteTemplate(w, tmpl+".html", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
