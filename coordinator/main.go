package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Item struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

func main() {
	port := "80"

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dummy := Item{"1", "Foo"}
			result, err := json.Marshal(dummy)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(result)
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
