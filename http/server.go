package http

import (
	"awesomeProject/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	port    int
	service *service.Service
}

func NewServer(port, concurrentParsers, dbBufferSize int) (*Server, error) {
	dbService, err := service.NewService(concurrentParsers, dbBufferSize)
	if err != nil {
		return nil, err
	}

	return &Server{
		port: port,
		service: dbService}, nil
}

func (s *Server) Start() error {
	svr := mux.NewRouter()

	svr.Path("/movies/load").
		Methods("POST").
		HandlerFunc(s.loadMovies)

	svr.Path("/movies/{id}").
		Methods("GET").
		HandlerFunc(getMovieByID)

	svr.Path("/movies/query/{query}").
		Methods("GET").
		HandlerFunc(getMovieByQuery)

	port := fmt.Sprintf(":%d", s.port)
	log.Printf("Listening on %s...", port)
	return http.ListenAndServe(port, svr)
}
