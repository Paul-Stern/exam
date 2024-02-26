package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Option struct {
	Id   int
	Text string
}

type Card struct {
	Id       int
	Question string
	Appendix []string
	Options  []Option
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Answer int

type Result struct {
	QuestionId string   `json:"questionId"`
	AnswerIds  []string `json:"answersIds"`
}

type TestResult struct {
	UserId  int      `json:"userId"`
	Results []Result `json:"results"`
}

type CardsResult []struct {
	QuestionId int   `json:"QuestionId"`
	AnswerIds  []int `json:"AnswerId"`
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
	Id            int       `json:"ID"`
	Task_text     string    `json:"TASK_TEXT"`
	Task_appendix []string  `json:"TASK_APPENDIX"`
	Answers       []restOpt `json:"ANSWERS"`
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
	urlPostRes   = "http://***REMOVED***:***REMOVED***/post"
)

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

var templates = template.Must(template.ParseFiles("test.html"))

func main() {
	http.HandleFunc("/test", makeHandler(viewHandler))
	http.HandleFunc("/login", signInHandler)
	http.HandleFunc("/json", jsonHandler)
	http.HandleFunc("/req", reqHandler)
	// Helps to test getting answers over post
	http.HandleFunc("/post", postHandler)
	log.Println("Server started. Listening to localhost:***REMOVED***")
	log.Fatal(http.ListenAndServe(":***REMOVED***", nil))
}

func newTest(u User, c []Card) Test {
	return Test{
		User:  getFullName(u),
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

func removeAppendixPrefix(ap []string) []string {
	if ap == nil {
		return ap
	}
	var result []string
	for _, s := range ap {
		_, j := utf8.DecodeRuneInString(s)
		// result[i] = s[(j * 3):]
		result = append(result, s[(j*3):])
	}
	return result
}

func newCardsResult(vals url.Values) TestResult {
	var tr TestResult
	tr.UserId = 1
	for k := range vals {
		var r Result
		// cr[k] = vals[k]
		qId, _ := strings.CutPrefix(k, "question_")
		var aIds []string
		for _, aId := range vals[k] {
			aIds = append(aIds, aId)
		}
		r.QuestionId = qId
		r.AnswerIds = aIds
		tr.Results = append(tr.Results, r)
	}
	return tr
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

func getPostJson(c Credentials, url string) (j message) {
	out, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("post error: %v", err)
	}
	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(out),
	)
	if err != nil {
		log.Fatalf("post error: %v", err)
	}
	// Read body
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("post error: %v", err)
	}
	return got
}
func sendPostJson(r CardsResult, url string) *http.Response {
	j, err := json.Marshal(r)
	if err != nil {
		log.Printf("sendPostJson error: %v", err)
	}
	res, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(j),
	)
	if err != nil {
		log.Printf("sendPostJson error: %v", err)
	}
	log.Println("success")

	return res
}

func getRestBlock(c Credentials) (rbs restBlocks) {
	got := getPostJson(c, urlQuestPost)
	err := json.Unmarshal(got, &rbs)
	if err != nil {
		log.Fatalf("getRestBlock error: %v", err)
	}
	return rbs
}

func getCards(rbs restBlocks) (cards []Card) {
	for _, block := range rbs {
		var c Card
		c.Id = block.Id
		c.Question = block.Task_text
		c.Appendix = removeAppendixPrefix(block.Task_appendix)
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
		log.Printf("%s", f)
		// cr := newCardResult(f)
		// w.Header().Add("Content-Type", "application/json")
		j, err := json.Marshal(newCardsResult(f))
		log.Printf("%s", j)
		if err != nil {
			log.Printf("post error: %v", err)
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", j)

		// log.Println(f)
		// http.Redirect(w, r, "/login", http.StatusFound)
	}
}
func postHandler(w http.ResponseWriter, r *http.Request) {
	b := getPostJson(getUserCreds(1), urlQuestPost)
	// b := r.Body
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", b)

}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "login.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Fatalf("ParseForm() err: %v", err)
		}

		c := readCreds(r.PostForm)
		u, err := getUserByEmail(c.Email)
		if err != nil {
			log.Printf("getUserByEmail() err: %v", err)
			http.Redirect(w, r, "/login", http.StatusFound)
		}

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
		for _, s := range sessions {
			u, err := getUserByEmail(s.email)
			if err != nil {
				log.Printf("Get current user error: %v", err)
			}
			c := getCards(getRestBlock(u.auth))
			t = newTest(u, c)
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
