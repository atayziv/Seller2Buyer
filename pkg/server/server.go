package server

import (
	"log"
	"net/http"

	"Practice/pkg/bid_request"
	"Practice/pkg/middleware"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Http   *http.Server
}

func NewServer() *Server {
	s := &Server{}
	s.Router = mux.NewRouter()
	s.mapRoutes()
	s.Http = &http.Server{
		Addr:    ":8080",
		Handler: s.Router,
	}
	return s
}

func (s *Server) mapRoutes() {
	s.Router.HandleFunc("/bid_request", bid_request.HandleBidRequest).Methods("POST")
	s.Router.Use(middleware.Validate)
}

func (s *Server) Start() {
	log.Println("Listening on port 8080...")
	if err := s.Http.ListenAndServe(); err != nil {
		log.Printf("Server error: %v", err)
	}
}
