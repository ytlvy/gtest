package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/ytlvy/gtest/gif"
)

var count int

type LServer struct {
	mu sync.Mutex
}

func NewLServer() *LServer {
	return &LServer{}
}

func (s *LServer) Run() {
	s.openLogFile("development.log")

	http.HandleFunc("/", s.handler)
	http.HandleFunc("/count", s.counter)
	log.Fatal(http.ListenAndServe(":8000", s.logRequest(http.DefaultServeMux)))
}

func (s *LServer) handler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	count++
	s.mu.Unlock()

	// fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	gif := new(gif.Maker)
	gif.MakeGif(w)
}

func (s *LServer) counter(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	fmt.Fprintf(w, "Count = %d\n", count)
	s.mu.Unlock()
}

func (s *LServer) logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func (s *LServer) openLogFile(path string) {
	if path != "" {
		lf, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			log.Fatal("open log file err :", err)
		}

		log.SetOutput(lf)
	}
}
