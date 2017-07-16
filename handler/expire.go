package handler

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func expireHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if err := store.Expire(ctx); err != nil {
		log.Errorf(ctx, "expireHandler: %T: %v", err, err)
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
