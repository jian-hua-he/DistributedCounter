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

type postItemHandler struct {
	Items []Item
}

func (h *postItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("In postItemHandler")
	switch r.Method {
	case http.MethodPost:
		var item Item
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.Items = append(h.Items, item)

		result, err := json.Marshal(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	default:
		log.Printf("Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

type getItemHandler struct {
	Items []Item
}

func (h *getItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("In getItemHandler")
	switch r.Method {
	case http.MethodGet:
		reg := regexp.MustCompile(`^\/items\/(.*)\/(count)$`)
		matchStr := reg.FindStringSubmatch(r.URL.Path)
		if len(matchStr) <= 1 {
			log.Printf("Regex not matched")
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
		log.Printf("Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func main() {
	items := []Item{}
	port := "80"

	http.Handle("/items/", &getItemHandler{Items: items})
	http.Handle("/items", &postItemHandler{Items: items})

	log.Printf("Start coordinator at %s port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
