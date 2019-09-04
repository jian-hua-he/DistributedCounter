package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

type Item struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

type postItemHandler struct{}

func (h *postItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var item Item
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := json.Marshal(item)
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
}

type countItemHandler struct{}

func (h *countItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
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
}

func main() {
	port := "80"

	http.Handle("/items", new(postItemHandler))

	http.Handle("/items/", new(countItemHandler))

	log.Printf("Start coordinator at %s port", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
