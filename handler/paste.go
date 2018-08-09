package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"syreclabs.com/dg/plain/storage"
)

const (
	contentKey = "plain"

	contentTypeHeader = "Content-Type"
	plainContentType  = "text/plain; charset=utf-8"
)

func pasteCreateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1024 * 1024); err != nil {
		writeError(w, "malformed request", http.StatusBadRequest)
		return
	}

	if len(r.PostForm) == 0 {
		writeError(w, "empty form data", http.StatusBadRequest)
		return
	}

	content := r.PostForm.Get(contentKey)
	if content == "" {
		writeError(w, "empty content", http.StatusBadRequest)
		return
	}

	ctx := appengine.NewContext(r)
	key, err := store.Add(ctx, content)
	if err != nil {
		log.Errorf(ctx, "pasteCreateHandler: %T: %v", err, err)
		writeError(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set(contentTypeHeader, plainContentType)
	w.Write([]byte(pasteUrl(r, key)))
}

func pasteShowHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get(":key")
	if key == "" {
		writeError(w, "missing paste key", http.StatusBadRequest)
		return
	}

	ctx := appengine.NewContext(r)
	content, err := store.Get(ctx, key)
	if err != nil {
		if err == storage.ErrNotFound {
			writeError(w, "not found", http.StatusNotFound)
		} else {
			log.Errorf(ctx, "pasteShowHandler: %T: %v", err, err)
			writeError(w, err.Error(), http.StatusServiceUnavailable)
		}
		return
	}

	w.Header().Set(contentTypeHeader, plainContentType)
	w.Write([]byte(content))
}

func pasteDeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get(":key")
	if key == "" {
		writeError(w, "missing paste key", http.StatusBadRequest)
		return
	}

	ctx := appengine.NewContext(r)
	if err := store.Delete(ctx, key); err != nil {
		if err == storage.ErrNotFound {
			writeError(w, "not found", http.StatusNotFound)
		} else {
			log.Errorf(ctx, "pasteDeleteHandler: %T: %v", err, err)
			writeError(w, err.Error(), http.StatusServiceUnavailable)
		}
		return
	}

	w.Header().Set(contentTypeHeader, plainContentType)
	w.Write([]byte("OK"))
}

func pasteUrl(r *http.Request, key string) string {
	url := url.URL{
		Scheme: r.URL.Scheme,
		Host:   r.Host,
	}
	if strings.HasSuffix(r.URL.Path, "/") {
		url.Path = r.URL.Path + key
	} else {
		url.Path = r.URL.Path + "/" + key
	}
	if url.Scheme == "" {
		url.Scheme = "http"
	}
	return url.String()
}

func writeError(w http.ResponseWriter, error string, code int) {
	http.Error(w, fmt.Sprintf("%d %s", code, error), code)
}
