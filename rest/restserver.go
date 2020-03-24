package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//Book book model
type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

//Server a simple restful server
type Server struct {
}

//NewServer return a new RestServer
func NewServer() *Server {
	return &Server{}
}

//Run start server
func (s *Server) Run() {
	r := httprouter.New()
	r.GET("/", s.homeHandler)

	//POST
	r.GET("/posts", s.postIndexHandler)
	r.POST("/posts", s.postCreateHandler)

	fmt.Println("starting server on :8080")
	http.ListenAndServe(":8080", r)
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// fmt.Fprintln(w, "home")
	book := Book{"Building Web Apps with Go", "Jeremy Saenz"}
	js, err := json.Marshal(book)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(js)
}

func (s *Server) postIndexHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "post index")
}

func (s *Server) postCreateHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "post create")
}

func (s *Server) postShowHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "post show")
}

func (s *Server) postUpdateHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "post update")
}

func (s *Server) postEditHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println(w, "post edit")
}
