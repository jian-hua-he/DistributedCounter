package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
)

type ItemService struct {
	Items []Item
}

type Item struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

type Count struct {
	Count int `json:"count"`
}

type ResponseStatus struct {
	Status string `json:"status"`
}

type postItemHandler struct {
	ItemService *ItemService
}

func (h *postItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("In postItemHandler")
	switch r.Method {
	case http.MethodPost:
		var items []Item
		err := json.NewDecoder(r.Body).Decode(&items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.ItemService.Items = append(h.ItemService.Items, items...)

		log.Printf("Current items: %+v", h.ItemService.Items)

		status := ResponseStatus{Status: "success"}
		result, _ := json.Marshal(status)

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

		itemMap := map[string]Item{}
		tenant := matchStr[1]
		for _, v := range h.ItemService.Items {
			if v.Tenant == tenant {
				itemMap[v.ID] = v
			}
		}
		count := Count{Count: len(itemMap)}

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

type healthHandler struct{}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success\n"))
}

func main() {
	service := ItemService{
		Items: []Item{},
	}
	port := "80"

	http.Handle("/items/", &getItemHandler{ItemService: &service})
	http.Handle("/items", &postItemHandler{ItemService: &service})
	http.Handle("/health", &healthHandler{})

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Start %s at %s port", hostname, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
