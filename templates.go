package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

const (
	templatesDir = "templates"
)

var (
	//go:embed templates/*
	files     embed.FS
	templates map[string]*template.Template
	// templates = template.Must(template.New("test").Funcs(funcMap).ParseFiles("templates/test.html"))
	// // go:embed templates/test.html
	funcMap = template.FuncMap{
		"getFullName": User.getFullName,
	}
)

// var funcMap = template.FuncMap{
// 	"getFullName": User.getFullName,
// }

func LoadTemplates() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tmpFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		return err
	}

	for _, tmpl := range tmpFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.New("test").Funcs(funcMap).ParseFS(files, templatesDir+"/"+tmpl.Name())
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}
	return nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, test *Test) {
	// err := templates[tmpl+".html"].ExecuteTemplate(w, tmpl+".html", t)
	t, ok := templates[tmpl+".html"]
	if !ok {
		log.Printf("template %s not found", tmpl+".html")
		return
	}

	if err := t.Execute(w, test); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}