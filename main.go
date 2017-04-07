package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	h := httpHandler{
		devMode: os.Getenv("MD_SLIDER_DEV_MODE") != "",
	}
	http.Handle("/", h)
	log.Print(http.ListenAndServe(":8000", nil))
}
