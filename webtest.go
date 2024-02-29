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
	email  string
	expiry time.Time
}

type message []byte

var sessions = map[string]session{}

var cfg Config

//go:embed static/login.html
var loginPage string

func main() {
	readConf(&cfg)
	err := LoadTemplates()
	if err != nil {
		log.Fatalf("LoadTemplates error: %v", err)
	}

	// http.Handle("/", http.FileServer(http.FS(loginFS)))

	http.HandleFunc("/test", makeHandler(viewHandler))
	http.HandleFunc("/login", signInHandler)
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
		log.Fatalf("getRestBlock error: %v", err)
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

		c := readCreds(r.PostForm)
		u, err := getUserByEmail(c.Email)
		if err != nil {
			log.Printf("getUserByEmail() err: %v", err)
			http.Redirect(w, r, "/login", http.StatusFound)
		}

		// Check if credentials are correct
		if c == u.Auth {
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

func getSession(u User) session {
	for uuid := range sessions {
		s := sessions[uuid]
		log.Printf("%v\n", s.email)
		return s
	}
	return session{}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, Test)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t Test
		for _, s := range sessions {
			u, err := getUserByEmail(s.email)
			if err != nil {
				log.Printf("Get current user error: %v", err)
			}
			c := getCards(getRestBlock(u.Auth))
			t = newTest(u, c)
		}
		fn(w, r, t)
	}
}
