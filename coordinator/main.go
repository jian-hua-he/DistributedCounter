package main

import (
	"log"
	"net/http"
)

type HostService struct {
	Hosts map[string]string
}

type ItemService struct {
	Items map[string]Item
}

type Item struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

type Count struct {
	Count int `json:"count"`
}

func main() {
	itemServ := ItemService{
		Items: map[string]Item{},
	}
	hostServ := HostService{
		Hosts: map[string]string{},
	}

	http.Handle("/items/", &getItemHandler{
		ItemService: &itemServ,
	})
	http.Handle("/items", &postItemHandler{
		ItemService: &itemServ,
		HostService: &hostServ,
	})
	http.Handle("/register", &regiesterHandler{
		HostService: &hostServ,
	})

	port := "80"
	log.Printf("Start coordinator at %s port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
