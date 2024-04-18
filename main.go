package main

import (
	_ "embed"
	"log"
	"net/http"
)

// Used to store value from LDFLAGS
var version string

func main() {
	err := LoadTemplates()
	if err != nil {
		log.Fatalf("LoadTemplates error: %v", err)
	}
	if version == "" {
		version = "dev"
	}

	// cert(255)
	log.Printf("Version: %s\n", version)
	testEmail()
	http.HandleFunc("/", root)
	http.HandleFunc("/login", signInHandler)
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/profiles", authenticate(profilesHandler))
	http.HandleFunc("/test", authenticate(testHandler))
	http.HandleFunc("/result", authenticate(resultHandler))
	http.HandleFunc("/logout", authenticate(logout))
	http.HandleFunc("/success", successHandler)
	// paths to the cert and the key
	log.Printf("Server started. Listening to %s:%s", cfg.Server.Addr, cfg.Server.Port)
	log.Fatal(http.ListenAndServeTLS(cfg.Server.Addr+":"+cfg.Server.Port, cfg.Server.Cert, cfg.Server.Key, nil))
}

func init() {
	readConf(&cfg)
}
