package main

import (
	"bytes"
	"encoding/json"
	"errors"
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
	log.Printf("%s %s", r.Method, r.URL.String())

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
	log.Printf("%s %s", r.Method, r.URL.String())

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
	log.Printf("%s %s", r.Method, r.URL.String())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success\n"))
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := registerHost(hostname); err != nil {
		log.Fatal(err)
	}

	service := ItemService{
		Items: []Item{},
	}
	http.Handle("/items/", &getItemHandler{ItemService: &service})
	http.Handle("/items", &postItemHandler{ItemService: &service})
	http.Handle("/health", &healthHandler{})

	port := "80"
	log.Printf("Start %s at %s port", hostname, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func registerHost(hostname string) error {
	data := []byte(hostname)
	resp, err := sendRequest("POST", "http://coordinator/register", bytes.NewBuffer(data))
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Register fail")
	}

	return nil
}

func sendRequest(method string, url string, buf *bytes.Buffer) (*http.Response, error) {
	var req *http.Request
	var err error

	// Pass buf directly will cause nil pointer error if buf is nil
	if buf == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, buf)
	}

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}
