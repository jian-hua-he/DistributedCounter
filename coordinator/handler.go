package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type postItemHandler struct {
	ItemService *ItemService
	HostService *HostService
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

type regiesterHandler struct {
	HostService *HostService
}

func (h *regiesterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("In registerHandler")

	switch r.Method {
	case http.MethodPost:
		// TODO: Check the token to ensure the request is from counter

		hostBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print("Error: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("fail\n"))
			return
		}

		hostname := string(hostBytes)
		url := fmt.Sprintf("http://%s/health", hostname)
		resp, err := sendRequest("GET", url, nil)
		if err != nil {
			log.Print("Error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if resp.StatusCode != http.StatusOK {
			log.Print("Error: health check fail")
			http.Error(w, "health check fail", http.StatusInternalServerError)
			return
		}

		h.HostService.Hosts[hostname] = hostname
		log.Printf("All hosts: %+v", h.HostService.Hosts)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success\n"))

	default:
		log.Printf("Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
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
