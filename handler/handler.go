package handler

import (
	"net/http"

	"github.com/dmgk/plain/storage"
	"github.com/gorilla/pat"
)

var store storage.Store

func New(s storage.Store) http.Handler {
	store = s

	r := pat.New()
	r.Get("/_cron/expire", expireHandler)
	r.Get("/{key}", pasteShowHandler)
	r.Delete("/{key}", pasteDeleteHandler)
	r.Post("/", pasteCreateHandler)
	r.Get("/", homeHandler)
	return r
}
