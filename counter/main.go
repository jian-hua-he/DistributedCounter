package main

import (
	"bytes"
	"encoding/json"
	"errors"
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
