package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
)

type ResponseStatus struct {
	Status string `json:"status"`
}

type postItemHandler struct {
	HostService *HostService
}

func (h *postItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodPost:
		resBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: Need to implement 2PC to handle error
		errors := []error{}
		var wg sync.WaitGroup
		for _, host := range h.HostService.Hosts {
			wg.Add(1)
			go func(host string) {
				url := fmt.Sprintf("http://%s/items", host)
				resp, err := sendRequest("POST", url, bytes.NewBuffer(resBodyBytes))
				if err != nil {
					log.Printf("Error: %s", err.Error())
					errors = append(errors, err)
				}
				defer func(resp *http.Response) {
					if resp != nil {
						resp.Body.Close()
					}
				}(resp)
				wg.Done()
			}(host)
		}
		wg.Wait()

		if len(errors) != 0 {
			status := ResponseStatus{Status: "fail"}
			result, _ := json.Marshal(status)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(result)
			return
		}

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

type getItemHandler struct{}

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

type registerHandler struct {
	HostService *HostService
}

func (h *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.String())

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
