package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Host: The Counter’s host
type Host struct {
	Name     string
	IsNew    bool
	Attempts int
}

// HostService: All Counters register in here
type HostService struct {
	Lock  sync.RWMutex
	Hosts map[string]Host
}

// Item: The data that save in the Counter
type Item struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

// Vote: For 2PC implementation
type Vote struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Data      []Item    `json:"data"`
}

// UpdateHost: Update hosts
// Golang map is not save to read and write
// Need to use the lock to prevent system collapse
func (hs *HostService) UpdateHost(h Host) {
	hs.Lock.Lock()
	defer hs.Lock.Unlock()

	hs.Hosts[h.Name] = h
}

// DeleteHost: Remove host
// Golang map is not save to read and write
// Need to use the lock to prevent system collapse
func (hs *HostService) DeleteHost(hostname string) {
	hs.Lock.Lock()
	defer hs.Lock.Unlock()

	delete(hs.Hosts, hostname)
}

// CheckHealth: Send the request to check the health in all Counters
// If check counting is up to maxAttempts
// Remove host from HostService
func (hs HostService) CheckHealth(maxAttempts int) {
	for _, host := range hs.Hosts {
		if host.Attempts >= maxAttempts {
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
