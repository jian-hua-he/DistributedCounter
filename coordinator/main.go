package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Host struct {
	Name     string
	IsNew    bool
	Attempts int
}

type HostService struct {
	Lock  sync.RWMutex
	Hosts map[string]Host
}

type Item struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

type Vote struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Data      []Item    `json:"data"`
}

func (hs *HostService) UpdateHost(h Host) {
	hs.Lock.Lock()
	defer hs.Lock.Unlock()

	hs.Hosts[h.Name] = h
}

func (hs *HostService) DeleteHost(hostname string) {
	hs.Lock.Lock()
	defer hs.Lock.Unlock()

	delete(hs.Hosts, hostname)
}

func (hs HostService) CheckHealth() {
	maxAttempts := 3

	for _, host := range hs.Hosts {
		if host.Attempts == maxAttempts {
			log.Printf("INFO: remove host %s", host.Name)
			hs.DeleteHost(host.Name)
			continue
		}

		url := fmt.Sprintf("http://%s/health", host.Name)
		resp, err := GET(url)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			host.Attempts += 1
			hs.UpdateHost(host)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			host.Attempts += 1
			hs.UpdateHost(host)
			continue
		}

		host.Attempts = 0
		hs.UpdateHost(host)
	}
}

func HealthCheck(hostServ *HostService) {
	// d := time.Duration(time.Minute * 1)
	d := time.Duration(time.Second * 30)
	t := time.NewTicker(d)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			log.Printf("INFO: start health check")
			log.Printf("INFO: all hosts %+v", hostServ.Hosts)
			hostServ.CheckHealth()
		}
	}
}

func main() {
	hostServ := HostService{
		Hosts: map[string]Host{},
	}

	go HealthCheck(&hostServ)

	http.Handle("/items/", &ItemCountHandler{})
	http.Handle("/items", &ItemHandler{
		HostService: &hostServ,
	})
	http.Handle("/register", &RegisterHandler{
		HostService: &hostServ,
	})
	http.Handle("/sync", &SyncHandler{
		HostService: &hostServ,
	})

	port := "80"
	log.Printf("INFO: start coordinator at %s port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
