package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

type ItemService struct {
	Items map[string]Item
}

type Item struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

type Count struct {
	Count int `json:"count"`
}

type postItemHandler struct {
	ItemService *ItemService
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

		h.ItemService.Items[item.ID] = item

		log.Printf("Current items: %+v", h.ItemService.Items)

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
	ItemService *ItemService
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

		tenant := matchStr[1]
		count := Count{Count: 0}
		for _, v := range h.ItemService.Items {
			if v.Tenant == tenant {
				count.Count += 1
			}
		}

		result, err := json.Marshal(count)
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
	service := ItemService{
		Items: map[string]Item{},
	}
	port := "80"

	http.Handle("/items/", &getItemHandler{ItemService: &service})
	http.Handle("/items", &postItemHandler{ItemService: &service})

	log.Printf("Start coordinator at %s port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
