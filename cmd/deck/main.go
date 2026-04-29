package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/nobe4/deck/internal/media"
	"github.com/nobe4/deck/internal/web"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	refreshMs := flag.Int("refresh", 5000, "refresh interval (ms)")
	debounceMs := flag.Int("debounce", 500, "debounce delay (ms)")
	tmplPath := flag.String("template", "", "path to custom HTML template")
	flag.Parse()

	ctrl, err := media.New()
	if err != nil {
		log.Fatalf("media: %v", err)
	}
	cfg := web.Config{
		RefreshMs:    *refreshMs,
		DebounceMs:   *debounceMs,
		TemplatePath: *tmplPath,
	}
	srv := web.New(ctrl, cfg)

	log.Printf("deck listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, srv))
}
