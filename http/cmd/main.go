package main

import (
	"awesomeProject/http"
	"log"
)

func main() {
	s, err := http.NewServer(5000, 15, 500)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
	log.Fatal(s.Start())
}