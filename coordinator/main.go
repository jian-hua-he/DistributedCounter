package main

import (
	"log"
	"net/http"
	"sync"
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
func main() {
	hostServ := HostService{
		Hosts: map[string]Host{},
	}

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
