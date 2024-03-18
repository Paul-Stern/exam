package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/kardianos/service"
)

// Used to store value from LDFLAGS
var version string

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// cert(255)
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
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	log.Println("Server stopped")
	return nil
}

func main() {
	readConf(&cfg)
	err := LoadTemplates()
	if err != nil {
		log.Fatalf("LoadTemplates error: %v", err)
	}
	if version == "" {
		version = "dev"
	}
	svcConfig := &service.Config{
		Name:        "WebtestService",
		DisplayName: "Webtest Service",
		Description: "Веб-сервер системы тестирования",
		Option: service.KeyValue{
			"UserService": true,
			"Interactive": true,
		},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Install()
	if err != nil {
		logger.Error(err)
	} else {
		log.Println("Service installed")
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}
