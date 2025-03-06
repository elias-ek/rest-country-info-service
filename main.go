package main

import (
	"log"
	"net/http"
	"time"

	"assignment-1/conf"
	"assignment-1/handlers"
)

func main() {

	//start timer for tracking uptime
	conf.StartTime = time.Now()

	//assign handlers to paths
	http.HandleFunc("/countryinfo/v1/info/", handlers.InfoHandler)
	http.HandleFunc("/countryinfo/v1/population/", handlers.PopHandler)
	http.HandleFunc("/countryinfo/v1/status/", handlers.StatHandler)

	port := ":8080"
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
