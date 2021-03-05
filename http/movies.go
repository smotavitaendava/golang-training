package http

import (
	"fmt"
	"net/http"
)

func (s *Server) loadMovies(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	body := req.Body
	defer body.Close()

	err := s.service.Process(body)
	var output string
	if err != nil {
		output = fmt.Sprintf(`{"success": %t, "error": "%s"}`, false, err)
	} else {
		output = fmt.Sprintf(`{"success": %t, "error": "%s"}`, true, "")
	}

	_, _ = res.Write([]byte(output))
}

func getMovieByID(res http.ResponseWriter, req *http.Request) {

}

func getMovieByQuery(res http.ResponseWriter, req *http.Request) {

}
