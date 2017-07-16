// +build appengine
package plain

import (
	"net/http"

	"github.com/enodata/plain/handler"
	"github.com/enodata/plain/storage"
)

func init() {
	store := storage.NewGDS()
	r := handler.New(store)
	http.Handle("/", r)
}
