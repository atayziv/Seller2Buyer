package server

import (
	"log"
	"net/http"

	"Practice/buyer_server/pkg/bid_response"

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
		Addr:    ":8081",
		Handler: s.Router,
	}
	return s
}

func (s *Server) mapRoutes() {
	s.Router.HandleFunc("/bid_response", bid_response.HandleBidRequest).Methods("POST")
}

func (s *Server) Start() {
	log.Println("Listening on port 8081...")
	if err := s.Http.ListenAndServe(); err != nil {
		log.Printf("Server error: %v", err)
	}
}
