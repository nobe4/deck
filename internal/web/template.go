package web

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"path/filepath"
)

//go:embed static/*
var staticFS embed.FS

func renderIndex(cfg Config) ([]byte, error) {
	var tmpl *template.Template
	var tmplName string
	var err error
	
	if cfg.TemplatePath != "" {
		tmpl, err = template.ParseFiles(cfg.TemplatePath)
		if err != nil {
			return nil, fmt.Errorf("parse template file: %w", err)
		}
		tmplName = filepath.Base(cfg.TemplatePath)
	} else {
		tmpl, err = template.ParseFS(staticFS, "static/*")
		if err != nil {
			return nil, fmt.Errorf("parse embedded templates: %w", err)
		}
		tmplName = "index.html"
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, tmplName, cfg); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	return buf.Bytes(), nil
}
