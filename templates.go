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

// type TemplateDataTypes interface {
// 	Profiles | ResultStore | Test
// }

// type TemplateData[DT TemplateDataTypes] struct {
// 	Data    DT
// 	Session session
// }

type TemplateData struct {
	Data    any
	Session session
}

var (
	//go:embed templates/*
	files     embed.FS
	templates map[string]*template.Template

	funcMap = template.FuncMap{
		"getFullName": User.getFullName,
		"inc": func(i int) int {
			return i + 1
		},
		"type": Card.Type,
	}
)

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

		t := template.New(tmpl.Name())

		t.Funcs(funcMap)

		pt, err := t.ParseFS(files, templatesDir+"/"+tmpl.Name())
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}
	return nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	// err := templates[tmpl+".html"].ExecuteTemplate(w, tmpl+".html", t)
	t, ok := templates[tmpl+".html"]
	if !ok {
		log.Printf("template %s not found", tmpl+".html")
		return
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
