package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/russross/blackfriday"
	"github.com/ytlvy/gtest/gif"
)

var count int

// LServer 一个简单 webserver
type LServer struct {
	mu sync.Mutex
}

// 包初始化函数
func init() {
	fmt.Println("Lserver init function called")
}

// NewLServer 返回一个 Lserver 实例
func NewLServer() *LServer {
	return &LServer{}
}

// Run start server
func (s *LServer) Run() {
	s.openLogFile("development.log")

	http.HandleFunc("/markdown", s.generateMarkdown)
	http.HandleFunc("/count", s.counter)

	http.Handle("/", http.FileServer(http.Dir("public")))

	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, s.logRequest(http.DefaultServeMux)))
}

//GenerateMarkdown 生成 markdown
func (s *LServer) generateMarkdown(w http.ResponseWriter, r *http.Request) {
	markdown := blackfriday.MarkdownCommon([]byte(r.FormValue("body")))
	w.Write(markdown)
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
