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
	"time"
)

type ResponseStatus struct {
	Status string `json:"status"`
}

type ItemHandler struct {
	HostService *HostService
}

type ItemCountHandler struct{}

type RegisterHandler struct {
	HostService *HostService
}

type SyncHandler struct {
	HostService *HostService
}

func (h *ItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodPost:

		// Validate item
		var items []Item
		if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
			log.Printf("ERROR: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Start 2PC
		// Vote phase
		now := time.Now()
		vote, _ := json.Marshal(Vote{
			ID:        GenID(now),
			Timestamp: now,
			Data:      items,
		})
		voteResults := make([]bool, 0)
		var wg sync.WaitGroup
		for _, host := range h.HostService.Hosts {
			wg.Add(1)
			go func(host Host) {
				url := fmt.Sprintf("http://%s/vote", host.Name)
				resp, err := POST(url, bytes.NewBuffer(vote))
				defer func(resp *http.Response) {
					if resp != nil {
						resp.Body.Close()
					}
					wg.Done()
				}(resp)
				if err != nil {
					log.Printf("ERROR: error occurred in vote phase %s", err.Error())
					voteResults = append(voteResults, false)
					return
				}

				if resp.StatusCode == http.StatusOK {
					voteResults = append(voteResults, true)
				} else {
					voteResults = append(voteResults, false)
				}
			}(host)
		}
		wg.Wait()

		// Check vote results
		rollback := false
		for _, v := range voteResults {
			if v == false {
				rollback = true
				break
			}
		}

		// Rollback phase
		if rollback {
			var wg sync.WaitGroup
			for _, host := range h.HostService.Hosts {
				wg.Add(1)
				go func(host Host) {
					url := fmt.Sprintf("http://%s/rollback", host.Name)
					resp, err := POST(url, bytes.NewBuffer(vote))
					defer func(resp *http.Response) {
						if resp != nil {
							resp.Body.Close()
						}
						wg.Done()
					}(resp)
					if err != nil {
						log.Printf("ERROR: error occurred in rollback phase. %s", err.Error())
						return
					}
				}(host)
			}
			wg.Wait()

			status := ResponseStatus{Status: "failed"}
			result, _ := json.Marshal(status)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(result)
			return
		}

		// Commit phase
		var wgCommit sync.WaitGroup
		for _, host := range h.HostService.Hosts {
			wgCommit.Add(1)
			go func(host Host) {
				url := fmt.Sprintf("http://%s/commit", host.Name)
				resp, err := POST(url, bytes.NewBuffer(vote))
				defer func(resp *http.Response) {
					if resp != nil {
						resp.Body.Close()
					}
					wgCommit.Done()
				}(resp)
				if err != nil {
					log.Printf("ERROR: error occurred in commit phase. %s", err.Error())
					return
				}
				if host.IsNew {
					host.IsNew = false
					h.HostService.UpdateHost(host)
				}
			}(host)
		}
		wgCommit.Wait()

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

		resp, err := GET("http://counter/" + matchStr[0])
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

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		host := Host{
			Name:     hostname,
			IsNew:    true,
			Attempts: 0,
		}
		h.HostService.UpdateHost(host)
		log.Printf("INFO: all hosts %+v", h.HostService.Hosts)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success\n"))

	default:
		log.Print("INFO: Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *SyncHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodGet:
		// TODO: Check the token to ensure the request is from counter

		for hostname, host := range h.HostService.Hosts {
			if host.IsNew {
				continue
			}

			url := fmt.Sprintf("http://%s/items", hostname)
			resp, err := GET(url)
			if err != nil {
				log.Printf("ERROR: error during sync. %s", err.Error())
				continue
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("ERROR: error during sync. response code is %s", resp.StatusCode)
				continue
			}

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("ERROR: error during sync. %s", err.Error())
				continue
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(bodyBytes)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))

	default:
		log.Print("INFO: Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func GET(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func POST(url string, buf *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}
