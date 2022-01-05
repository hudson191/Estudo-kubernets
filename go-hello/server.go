package main

import (
	"fmt"
	"net/http"
	"time"
)

var startedAt = time.Now()

func main() {
	http.HandleFunc("/", Hello)
	http.HandleFunc("/healthz", Health)
	http.ListenAndServe(":80", nil)
}

func Health(w http.ResponseWriter, r *http.Request) {
	duration := time.Since(startedAt)

	if duration.Seconds() > 25 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Duration: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("OK %v", duration.Seconds())))
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Batatinha Frita</h1>"))
}
