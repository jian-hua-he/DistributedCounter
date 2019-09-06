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
	"time"
)

type Transaction struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Data      []Item    `json:"data"`
}

type ItemService struct {
	Items        []Item
	Transactions []Transaction
}

func (s *ItemService) TransAppend(t Transaction) {
	s.Transactions = append(s.Transactions, t)
}

func (s *ItemService) DoTrans(id string) {
	for i, t := range s.Transactions {
		if t.ID != id {
			continue
		}

		s.Items = append(s.Items, t.Data...)
		s.Transactions = append(s.Transactions[:i], s.Transactions[i+1:]...)
		break
	}
}

func (s *ItemService) AbortTrans(id string) {
	for i, t := range s.Transactions {
		if t.ID != id {
			continue
		}

		s.Transactions = append(s.Transactions[:i], s.Transactions[i+1:]...)
		break
	}
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
		log.Printf("ERROR: error occurr during register. %s", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Register status code %v", resp.StatusCode)
		err := errors.New(msg)
		log.Printf("ERROR: error occurr during register. %s", err.Error())
		return err
	}

	return nil
}

func SyncItems() ([]Item, error) {
	log.Print("INFO: start to sync process")

	resp, err := GET("http://coordinator/sync")
	defer resp.Body.Close()
	if err != nil {
		log.Printf("ERROR: error occured during sync. %s", err.Error())
		return []Item{}, err
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Sync status code %v", resp.StatusCode)
		err := errors.New(msg)
		log.Printf("ERROR: error occured during sync. %s", err.Error())
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
	if err := RegisterHost(hostname); err != nil {
		log.Fatal("ERROR: " + err.Error())
	}
	items, err := SyncItems()
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
