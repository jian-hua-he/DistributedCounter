package main

import (
	"log"
	"net/http"
)

type HostService struct {
	// string is the hostname
	// bool is represent new counter
	Hosts map[string]bool
}

func main() {
	hostServ := HostService{
		Hosts: map[string]bool{},
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
