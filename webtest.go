package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Answer int

type session struct {
	user    User
	profile TestProfile
	expiry  time.Time
	start   time.Time
	result  ResultStore
}

type DataTypes interface {
	User | Tasks | TestProfile | []TestProfile | AvailableTestProfiles | Profiles | ResultStore
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

var version string

func main() {
	readConf(&cfg)
	err := LoadTemplates()
	if err != nil {
		log.Fatalf("LoadTemplates error: %v", err)
	}

	http.HandleFunc("/login", signInHandler)
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/profiles", profilesHandler)
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/success", successHandler)
	// Helps to test getting answers over post
	log.Printf("Version: %s\n", version)
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
	// Create json
	out, err := json.Marshal(v)
	if err != nil {
		log.Printf("post error: %v", err)
	}

	// Post json and get reponse
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

func testHandler(w http.ResponseWriter, r *http.Request) {
	// get session id from cookie
	sesCookie, _ := r.Cookie("gosesid")
	// get session data
	ses := sessions[sesCookie.Value]
	switch r.Method {
	case "GET":
		var t Test
		// get user from session
		profid := strconv.Itoa(ses.profile.Id)
		// get tasks
		tasks, _ := getTasks(profid)
		c := getCards(tasks)
		// Create test
		t = newTest(ses.user, c, ses.profile)
		ses.start = t.Time.Start
		http.SetCookie(w, &http.Cookie{
			Name:  "testing_start",
			Value: t.Time.Start.Format(time.RFC3339),
		})
		// http.SetCookie(w, profid)
		renderTemplate(w, "test", &t)
	case "POST":
		// Get data from form
		if err := r.ParseForm(); err != nil {
			log.Printf("ParseForm() err: %v", err)
			return
		}
		f := r.PostForm
		// Put data to Test Result
		tr, err := newTestResult(f)
		if err != nil {
			log.Printf("test error: %v", err)
		}
		testStart, _ := r.Cookie("testing_start")
		tr.Time.Start, err = time.Parse(time.RFC3339, testStart.Value)
		if err != nil {
			log.Printf("test error: %v", err)
		}
		url := baseUrl(cfg) + "/" + "tests"

		// Post test results and get response
		resp, err := post(tr, url)
		if err != nil {
			log.Printf("test error: %v", err)
		}
		// Parse stored result id
		result, _ := read[ResultStore](resp)
		// Save stored result id to session
		ses.result.Id = result.Data.Id
		// save session
		sessions[sesCookie.Value] = ses

		http.SetCookie(w, sesCookie)
		http.Redirect(w, r, "/result", http.StatusFound)
	}
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
		cookie.Path = "/"
		http.SetCookie(w, &cookie)
		log.Println("Login: success. Redirecting...")
		http.Redirect(w, r, "/profiles", http.StatusFound)
	}
}

// TODO: Reimplement this function
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
		http.Redirect(w, r, "/success", http.StatusFound)

	}
}

func profilesHandler(w http.ResponseWriter, r *http.Request) {
	// Get session id from cookie
	sesCookie, _ := r.Cookie("gosesid")
	ses := sessions[sesCookie.Value]

	switch r.Method {
	case http.MethodGet:
		var profiles Profiles
		// get tasks
		profiles, _ = getProfiles()
		renderTemplate(w, "profiles", &profiles)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			log.Printf("ParseForm() err: %v", err)
			return
		}
		log.Printf("form: %+v\n", r.Form)

		pids := r.FormValue("TASK_PROFILE_ID")
		pid, _ := strconv.Atoi(pids)
		ptext := r.FormValue("TASK_PROFILE_TEXT_" + pids)
		profile := TestProfile{
			Id:   pid,
			Text: ptext,
		}

		ses.profile = profile
		// Save session
		sessions[sesCookie.Value] = ses
		http.Redirect(w, r, "/test", http.StatusFound)

	}

}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	// get session id from cookie
	sesCookie, _ := r.Cookie("gosesid")
	// get session data
	ses := sessions[sesCookie.Value]
	result, _ := getResult(ses.result.Id)
	// test scenario
	// result, _ := getResult(255)
	renderTemplate(w, "result", result)
}

func successHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, successPage)
}

func getTasks(profileId string) (tasks Tasks, err error) {
	// Build url
	url := strings.Join([]string{
		baseUrl(cfg),
		"profiles",
		profileId,
		"tasks",
	}, "/")
	// Send request and get response
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	// Parse response to message
	m, err := read[Tasks](resp)
	if err != nil {
		return
	}
	// Save message data to tasks
	tasks = m.Data
	return
}

func getProfiles() (profiles Profiles, err error) {
	url := baseUrl(cfg) + "/profiles"
	// Send request
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	// Parse response
	m, err := read[Profiles](resp)
	if err != nil {
		return
	}
	// Save message data to profiles
	profiles = m.Data
	return
}

func getResult(id int) (result ResultStore, err error) {
	ids := strconv.Itoa(id)
	url := baseUrl(cfg) + "/tests/" + ids
	// Send request
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	// Parse response
	m, err := read[ResultStore](resp)
	if err != nil {
		return
	}
	result = m.Data
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
