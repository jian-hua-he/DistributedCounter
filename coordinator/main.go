package main

import (
	"log"
	"net/http"
)

type HostService struct {
	Hosts map[string]string
}

func main() {
	hostServ := HostService{
		Hosts: map[string]string{},
	}

	http.Handle("/items/", &getItemHandler{})
	http.Handle("/items", &postItemHandler{
		HostService: &hostServ,
	})
	http.Handle("/register", &registerHandler{
		HostService: &hostServ,
	})

	port := "80"
	log.Printf("Start coordinator at %s port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
