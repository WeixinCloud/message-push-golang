package main

import (
	"log"
	"message-push-golang/service"
	"net/http"
)

func main() {
	http.HandleFunc("/", service.IndexHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
