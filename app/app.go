package app

import (
	"net/http"

	"github.com/dmgk/plain/handler"
	"github.com/dmgk/plain/storage"
)

func init() {
	store := storage.NewGDS()
	r := handler.New(store)
	http.Handle("/", r)
}
