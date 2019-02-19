package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dmgk/plain/handler"
	"github.com/dmgk/plain/storage"
)

func main() {
	store := storage.NewGDS()
	r := handler.New(store)

	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
