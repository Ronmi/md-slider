package main

import "net/http"

func main() {
	http.HandleFunc("/", httpHandler)
	http.ListenAndServe(":8000", nil)
}
