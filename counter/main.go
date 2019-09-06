package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ItemService struct {
	Items []Item
}

type Item struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

type Count struct {
	Count int `json:"count"`
}

func RegisterHost(hostname string) error {
	data := []byte(hostname)
	resp, err := POST("http://coordinator/register", bytes.NewBuffer(data))
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Register fail")
	}

	return nil
}

func SyncItems() ([]Item, error) {
	resp, err := GET("http://coordinator/sync")
	defer resp.Body.Close()
	if err != nil {
		return []Item{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []Item{}, errors.New("Sync fail")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Item{}, err
	}

	var items []Item
	if err := json.Unmarshal(bodyBytes, &items); err != nil {
		return []Item{}, err
	}

	return items, nil
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("ERROR: " + err.Error())
	}
	if err := RegisterHost(hostname); err != nil {
		log.Fatal("ERROR: " + err.Error())
	}
	items, err := SyncItems()
	if err != nil {
		log.Fatal("ERROR: " + err.Error())
	}
	service := ItemService{
		Items: items,
	}

	http.Handle("/items/", &ItemCountHandler{ItemService: &service})
	http.Handle("/items", &ItemHandler{ItemService: &service})
	http.Handle("/vote", &VoteHandler{})
	http.Handle("/commit", &CommitHandler{})
	http.Handle("/rollback", &RollbackHandler{})
	http.Handle("/health", &HealthHandler{})

	port := "80"
	log.Printf("INFO: start %s at %s port", hostname, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
