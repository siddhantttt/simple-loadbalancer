package main

import (
	"log"
	"net/http"
)

func main() {
	lb := NewLoadBalancer("localhost:8080", "localhost:8081")
	go lb.StartHealthcheck()

	mux := http.NewServeMux()
	mux.Handle("/*", lb)

	err := http.ListenAndServe("localhost:8000", mux)
	if err != nil {
		log.Fatal("Something went wrong")
	}
}
