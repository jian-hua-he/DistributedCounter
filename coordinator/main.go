package main

import (
	"bytes"
	"io/ioutil"
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
		resBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := sendRequest("POST", "http://counter/items", bytes.NewBuffer(resBodyBytes))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bodyBytes)
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

		resp, err := sendRequest("GET", "http://counter/"+matchStr[0], nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bodyBytes)
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
