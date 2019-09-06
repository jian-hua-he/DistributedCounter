package main

import (
	"log"
	"net/http"
	"time"
)

// HealthCheck: Do health check in every 10 second
func HealthCheck(hostServ *HostService) {
	maxAttempts := 3
	d := time.Duration(time.Second * 10)
	t := time.NewTicker(d)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			log.Printf("INFO: start health check")
			log.Printf("INFO: all hosts %+v", hostServ.Hosts)
			hostServ.CheckHealth(maxAttempts)
		}
	}
}

func main() {
	hostServ := HostService{
		Hosts: map[string]Host{},
	}

	go HealthCheck(&hostServ)

	http.Handle("/items/", &ItemCountHandler{})
	http.Handle("/items", &ItemHandler{
		HostService: &hostServ,
	})
	http.Handle("/register", &RegisterHandler{
		HostService: &hostServ,
	})

	port := "80"
	log.Printf("INFO: start coordinator at %s port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
