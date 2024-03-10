package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
)

type Answer int

type session struct {
	user   User
	expiry time.Time
}

type DataTypes interface {
	User | Tasks | TestProfile | []TestProfile | AvailableTestProfiles | Profiles
}

type Message[D DataTypes] struct {
	Error struct {
		Code int    `json:"CODE"`
		Text string `json:"TEXT"`
	} `json:"ERROR"`
	Data D `json:"DATA"`
}

var sessions = map[string]session{}

var cfg Config

//go:embed static/login.html
var loginPage string

//go:embed static/signup.html
var signupPage string

//go:embed static/successfulRegistration.html
var successPage string

func main() {
	readConf(&cfg)
	err := LoadTemplates()
	if err != nil {
		log.Fatalf("LoadTemplates error: %v", err)
	}

	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/login", signInHandler)
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/profiles", profilesHandler)
	http.HandleFunc("/success", successHandler)
	http.HandleFunc("/json", jsonHandler)
	// Helps to test getting answers over post
	http.HandleFunc("/post", postHandler)
	log.Printf("Server started. Listening to localhost%s", ":"+cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, nil))
}
func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
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

func post(v any, url string) (*http.Response, error) {
	out, err := json.Marshal(v)
	if err != nil {
		log.Printf("post error: %v", err)
	}

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(out),
	)
	if err != nil {
		log.Printf("post error: %v", err)
	}
	return resp, err
}

func getPostJson(c Credentials, url string) []byte {
	out, err := json.Marshal(c)
	if err != nil {
		log.Printf("post error: %v", err)
	}
	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(out),
	)
	if err != nil {
		log.Printf("post error: %v", err)
	}
	// Read body
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("post error: %v", err)
	}
	return got
}
func sendPostJson(r TestResult, url string) *http.Response {
	j, err := json.Marshal(r)
	if err != nil {
		log.Printf("sendPostJson error: %v", err)
	}
	if !json.Valid(j) {
		err = errors.New("json is not valid")
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

func getRestBlock(c Credentials) (tasks Tasks) {
	got := getPostJson(c, getQuestionUrl(cfg))
	err := json.Unmarshal(got, &tasks)
	if err != nil {
		log.Printf("getRestBlock error: %v", err)
	}
	return tasks
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(users[1].Auth)
	if err != nil {
		log.Printf("Json error: %v", err)
	}
	w.Header().Add("Content-type", "application/json")
	fmt.Fprintf(w, "%s", data)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// var t Test
		// get user from cookie
		// cookie, _ := r.Cookie("Session id")
		// u := sessions[cookie.Value].user
		// // get tasks
		// tasks, _ := getTasks(u)
		// c := getCards(tasks)
		// t = newTest(u, c)
		// renderTemplate(w, "test", &t)
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Printf("ParseForm() err: %v", err)
			return
		}
		f := r.PostForm
		log.Printf("%s", f)
		tr, err := newTestResult(f)
		if err != nil {
			log.Printf("post error: %v", err)
		}
		log.Printf("%v", tr)
		r := sendPostJson(tr, getSaveUrl(cfg))
		got, err := io.ReadAll(r.Body)

		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			log.Fatalf("post error: %v", err)
		}
		fmt.Fprintf(w, "%s", got)
	}
}
func postHandler(w http.ResponseWriter, r *http.Request) {
	b := getPostJson(getUserCreds(1), getQuestionUrl(cfg))
	// b := r.Body
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", b)

}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprint(w, loginPage)
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Fatalf("ParseForm() err: %v", err)
		}

		// Read credentials from form
		c := readCreds(r.PostForm)
		// Get user by from REST server
		u, err := getUser(c.Email)
		if err != nil {
			return
		}
		// Check password
		if c.Password != u.Auth.Password {
			log.Println("login/password not correct")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		sessionToken := uuid.NewString()
		expiresAt := time.Now().Add(2 * time.Hour)

		// save data to session
		sessions[sessionToken] = session{
			user:   u,
			expiry: expiresAt,
		}
		cookie := http.Cookie{}
		cookie.Name = "gosesid"
		cookie.Value = sessionToken
		cookie.Path = "/profiles"
		http.SetCookie(w, &cookie)
		log.Println("Login: success. Redirecting...")
		http.Redirect(w, r, "/profiles", http.StatusFound)
	}
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprint(w, signupPage)
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Printf("ParseForm() err: %v", err)
			return
		}

		u := newUser(
			r.Form["name"][0],
			r.Form["middlename"][0],
			r.Form["surname"][0],
			r.Form["email"][0],
			r.Form["password"][0])
		// fmt.Fprintf(w, "%v", u)

		// Send registration data to REST server
		resp, _ := post(u, getRegister(cfg))
		// Parse response with user data from REST server
		read[User](resp)
		// j, _ := json.Marshal(m)
		// w.Header().Add("Content-Type", "application/json")
		// fmt.Fprintf(w, "%s", j)
		http.Redirect(w, r, "/success", http.StatusFound)

	}
}

func profilesHandler(w http.ResponseWriter, r *http.Request) {
	// get user from cookie
	// cookie, _ := r.Cookie("gosesid")
	// u := sessions[cookie.Value].user

	switch r.Method {
	case http.MethodGet:
		// get tasks
		profiles, _ := getProfiles()
		renderTemplate(w, "profiles", &profiles)
	case http.MethodPost:

	}

}

func successHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, successPage)
}

/*
	func getTasks(u User) (tasks Tasks, err error) {
		resp, err := post(u, getQuestionUrl(cfg))
		if err != nil {
			return
		}
		m, err := read[Tasks](resp)
		if err != nil {
			return
		}
		tasks = m.Data
		return
	}
*/

func getTasks(prof AvailableTestProfiles) (tasks Tasks, err error) {
	// resp, err := post(prof, "")
	return
}

func getProfiles() (profiles Profiles, err error) {
	url := baseUrl(cfg) + "/profiles"
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	m, err := read[Profiles](resp)
	if err != nil {
		return
	}
	profiles = m.Data
	return
}

func readCreds(f url.Values) Credentials {
	return Credentials{
		Email:    f["email"][0],
		Password: f["password"][0],
	}
}

func read[DT DataTypes](r *http.Response) (m Message[DT], err error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return
	}
	if m.Error.Code == 0 {
		err = nil
	}
	return
}
