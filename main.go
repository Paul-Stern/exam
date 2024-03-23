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
	testEmail()
	http.HandleFunc("/login", signInHandler)
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/profiles", profilesHandler)
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/success", successHandler)
	// Helps to test getting answers over post
	log.Printf("Version: %s\n", version)
	// paths to the cert and the key
	log.Printf("Server started. Listening to %s:%s", cfg.Server.Addr, cfg.Server.Port)
	log.Fatal(http.ListenAndServeTLS(cfg.Server.Addr+":"+cfg.Server.Port, cfg.Server.Cert, cfg.Server.Key, nil))
}

func init() {
	readConf(&cfg)
}
