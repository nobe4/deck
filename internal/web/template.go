package web

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"path/filepath"
)

//go:embed static/*
var staticFS embed.FS

func renderIndex(cfg Config) []byte {
	var tmpl *template.Template
	var tmplName string
	if cfg.TemplatePath != "" {
		tmpl = template.Must(template.ParseFiles(cfg.TemplatePath))
		tmplName = filepath.Base(cfg.TemplatePath)
	} else {
		tmpl = template.Must(template.ParseFS(staticFS, "static/*"))
		tmplName = "index.html"
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, tmplName, cfg); err != nil {
		log.Fatalf("render index: %v", err)
	}

	return buf.Bytes()
}
