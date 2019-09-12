package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

// ItemHandler: Handle http://{hostname}/items request
// Return all items for sync data usage
type ItemHandler struct {
	ItemService *ItemService
}

// ItemCountHandler: Handle http://{hostname}/items/{tenant}/count request
// Return the count of tenant
type ItemCountHandler struct {
	ItemService *ItemService
}

// VoteHandler: Handle http://{hostname}/vote request
// For first phase of 2PC
// Save the transaction to ItemService
type VoteHandler struct {
	ItemService *ItemService
}

// CommitHandler: Handle http://{hostname}/commit request
// For second phase (commit) of 2PC
// Put data from the transaction to items in ItemService
// After that, remove the transaction
type CommitHandler struct {
	ItemService *ItemService
}

// RollbackHandler: Handle http://{hostname}/rollback request
// For second phase (rollback) of 2PC
// Cancel and remove the transaction
type RollbackHandler struct {
	ItemService *ItemService
}

// HealthHandler: Handle http://{hostname}/health request
// For health check
type HealthHandler struct{}

func (h *ItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodGet:
		items, _ := json.Marshal(h.ItemService.Items)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(items)

	default:
		log.Printf("INFO: Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *ItemCountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodGet:
		// Check the url
		reg := regexp.MustCompile(`^\/items\/(.*)\/(count)$`)
		matchStr := reg.FindStringSubmatch(r.URL.Path)
		if len(matchStr) <= 1 {
			log.Print("INFO: Regex not matched")
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Calculate tenant count
		itemMap := map[string]Item{}
		tenant := matchStr[1]
		for _, v := range h.ItemService.Items {
			if v.Tenant == tenant {
				itemMap[v.ID] = v
			}
		}
		count := Count{Count: len(itemMap)}

		// Transfer to json string
		result, err := json.Marshal(count)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)

	default:
		log.Print("INFO: Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *VoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodPost:
		var tran Transaction
		err := json.NewDecoder(r.Body).Decode(&tran)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.ItemService.TransAppend(tran)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success\n"))

	default:
		log.Print("INFO: Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}

}

func (h *CommitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodPost:
		var tran Transaction
		err := json.NewDecoder(r.Body).Decode(&tran)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.ItemService.DoTrans(tran.ID)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success\n"))

	default:
		log.Print("INFO: Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *RollbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodPost:
		var tran Transaction
		err := json.NewDecoder(r.Body).Decode(&tran)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.ItemService.AbortTrans(tran.ID)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success\n"))

	default:
		log.Print("INFO: Unaccept method")
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s", r.Method, r.URL.String())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success\n"))
}

// GET: Send the request with GET method
func GET(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

// POST: Send the request with POST method
func POST(url string, buf *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}
