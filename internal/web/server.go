// Package web provides HTTP server and API handlers for media control.
package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/nobe4/deck/internal/media"
)

type Config struct {
	RefreshMs    int
	DebounceMs   int
	TemplatePath string
}

type Server struct {
	ctrl  media.Controller
	mux   *http.ServeMux
	index []byte
}

func New(ctrl media.Controller, cfg Config) (*Server, error) {
	index, err := renderIndex(cfg)
	if err != nil {
		return nil, err
	}

	s := &Server{
		ctrl:  ctrl,
		mux:   http.NewServeMux(),
		index: index,
	}

	s.mux.HandleFunc("GET /", s.serveIndex)
	s.mux.HandleFunc("GET /api/{field}", s.getField)
	s.mux.HandleFunc("POST /api/{field}", s.setField)

	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write(s.index); err != nil {
		log.Printf("write index: %v", err)
	}
}

func (s *Server) getField(w http.ResponseWriter, r *http.Request) {
	field := r.PathValue("field")

	val, err := s.get(field)
	if err != nil {
		log.Printf("get %s: %v", field, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{field: val}); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func (s *Server) setField(w http.ResponseWriter, r *http.Request) {
	field := r.PathValue("field")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.set(field, body); err != nil {
		log.Printf("set %s: %v", field, err)
		status := http.StatusInternalServerError
		if _, ok := err.(*json.SyntaxError); ok {
			status = http.StatusBadRequest
		} else if _, ok := err.(*json.UnmarshalTypeError); ok {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
