package main

import "net/http"

func main() {
	http.HandleFunc("/200",func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Status", "200")
	})
	http.HandleFunc("/500",func(w http.ResponseWriter, r *http.Request) {
		// NOTE: you can't set status code in header, it's the first line in header
		// and Status is just a normal header ...
		w.Header().Set("Status", "500")
	})
	http.HandleFunc("/5002",func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})
	http.ListenAndServe(":8000", http.DefaultServeMux)
}
