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

func New(ctrl media.Controller, cfg Config) *Server {
	s := &Server{
		ctrl:  ctrl,
		mux:   http.NewServeMux(),
		index: renderIndex(cfg),
	}

	s.mux.HandleFunc("GET /", s.serveIndex)
	s.mux.HandleFunc("GET /api/{field}", s.getField)
	s.mux.HandleFunc("POST /api/{field}", s.setField)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(s.index)
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
	_ = json.NewEncoder(w).Encode(map[string]any{field: val})
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
