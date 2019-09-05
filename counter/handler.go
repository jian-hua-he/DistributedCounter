package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

type ResponseStatus struct {
	Status string `json:"status"`
}

type ItemHandler struct {
	ItemService *ItemService
}

type ItemCountHandler struct {
	ItemService *ItemService
}

type HealthHandler struct{}

func (h *ItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodGet:
		items, _ := json.Marshal(h.ItemService.Items)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(items)

	case http.MethodPost:
		var items []Item
		err := json.NewDecoder(r.Body).Decode(&items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.ItemService.Items = append(h.ItemService.Items, items...)

		log.Printf("INFO: current items %+v", h.ItemService.Items)

		status := ResponseStatus{Status: "success"}
		result, _ := json.Marshal(status)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)

	default:
		log.Printf("INFO: Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *ItemCountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodGet:
		reg := regexp.MustCompile(`^\/items\/(.*)\/(count)$`)
		matchStr := reg.FindStringSubmatch(r.URL.Path)
		if len(matchStr) <= 1 {
			log.Print("INFO: Regex not matched")
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
		log.Print("INFO: Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success\n"))
}

func GET(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func POST(url string, buf *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}