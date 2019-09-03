package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := "80"

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Println("Get counter")
		case http.MethodPost:
			fmt.Println("Post items")
		default:
			fmt.Println("Invalid method")
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Test")
	})

	log.Printf("Start coordinator at %s port", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
