package app // import "syreclabs.com/dg/plain/app"

import (
	"net/http"

	"syreclabs.com/dg/plain/handler"
	"syreclabs.com/dg/plain/storage"
)

func init() {
	store := storage.NewGDS()
	r := handler.New(store)
	http.Handle("/", r)
}
