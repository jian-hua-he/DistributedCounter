package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// RegisterHost: Send the request to coordinator and get the initial items
func RegisterHost(hostname string) ([]Item, error) {
	data := []byte(hostname)
	resp, err := POST("http://coordinator/register", bytes.NewBuffer(data))
	defer func(resp *http.Response) {
		if resp != nil {
			resp.Body.Close()
		}
	}(resp)
	if err != nil {
		log.Printf("ERROR: error occurr during register. %s", err.Error())
		return []Item{}, err
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Register status code %v", resp.StatusCode)
		err := errors.New(msg)
		log.Printf("ERROR: error occurr during register. %s", err.Error())
		return []Item{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: error occured during sync. %s", err.Error())
		return []Item{}, err
	}

	var items []Item
	if err := json.Unmarshal(bodyBytes, &items); err != nil {
		log.Printf("ERROR: error occured during sync. %s", err.Error())
		return []Item{}, err
	}

	return items, nil
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("ERROR: " + err.Error())
	}

	items, err := RegisterHost(hostname)
	if err != nil {
		log.Fatal("ERROR: " + err.Error())
	}

	service := ItemService{
		Items:        items,
		Transactions: []Transaction{},
	}

	http.Handle("/items/", &ItemCountHandler{ItemService: &service})
	http.Handle("/items", &ItemHandler{ItemService: &service})
	http.Handle("/vote", &VoteHandler{ItemService: &service})
	http.Handle("/commit", &CommitHandler{ItemService: &service})
	http.Handle("/rollback", &RollbackHandler{ItemService: &service})
	http.Handle("/health", &HealthHandler{})

	port := "80"
	log.Printf("INFO: start %s at %s port", hostname, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
