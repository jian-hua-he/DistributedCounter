package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
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
			w.WriteHeader(http.StatusOK)
			w.Write(result)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})

	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			reg := regexp.MustCompile(`^\/items\/(.*)\/(count)$`)
			matchStr := reg.FindStringSubmatch(r.URL.Path)
			if len(matchStr) <= 1 {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			dummy := Item{
				ID:     matchStr[1],
				Tenant: "",
			}
			result, err := json.Marshal(dummy)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(result)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Test")
	})

	log.Printf("Start coordinator at %s port", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
