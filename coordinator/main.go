package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Test")
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}
