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

type getItemHandler struct{}

type registerHandler struct {
	HostService *HostService
}

func (h *postItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

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
		for host, _ := range h.HostService.Hosts {
			wg.Add(1)
			go func(host string) {
				url := fmt.Sprintf("http://%s/items", host)
				resp, err := sendRequest("POST", url, bytes.NewBuffer(resBodyBytes))
				if err != nil {
					log.Printf("ERROR: %s", err.Error())
					errors = append(errors, err)
				}
				defer func(resp *http.Response) {
					if resp != nil {
						h.HostService.Hosts[host] = false
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
		log.Print("INFO: Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *getItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

		resp, err := sendRequest("GET", "http://counter/"+matchStr[0], nil)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bodyBytes)

	default:
		log.Print("INFO: Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodPost:
		// TODO: Check the token to ensure the request is from counter

		hostBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print("ERROR: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("fail\n"))
			return
		}
		hostname := string(hostBytes)
		h.HostService.Hosts[hostname] = true
		log.Printf("INFO: all hosts %+v", h.HostService.Hosts)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success\n"))

	default:
		log.Print("INFO: Unaccept method")
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
