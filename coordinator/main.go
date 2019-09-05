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

	http.Handle("/items/", &getItemHandler{})
	http.Handle("/items", &postItemHandler{
		HostService: &hostServ,
	})
	http.Handle("/register", &registerHandler{
		HostService: &hostServ,
	})

	port := "80"
	log.Printf("INFO: Start coordinator at %s port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
