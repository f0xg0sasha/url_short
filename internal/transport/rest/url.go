package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/f0xg0sasha/url_short/internal/domain"
	"github.com/f0xg0sasha/url_short/internal/storage"
	"github.com/gorilla/mux"
)

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alias := vars["alias"]

	OriginalURL, err := h.urlService.Fetch(context.Background(), alias)
	if err != nil {
		if err == storage.ErrURLNotFound {
			http.Error(w, "error", http.StatusNotFound)
			h.log.Error(err)
			return
		}

		http.Error(w, "error", http.StatusBadRequest)
		h.log.Error(err)
		return
	}

	h.log.Info(fmt.Sprintf("succses redirect [alias]: (%s)", alias))
	http.Redirect(w, r, OriginalURL, http.StatusMovedPermanently)
}

func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	url := &domain.RequestURL{}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		h.log.Error(err)
		return
	}

	err = json.Unmarshal(reqBytes, url)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		h.log.Error(err)
		return
	}

	id, err := h.urlService.Create(context.Background(), url.URL, url.Alias)
	if err != nil {
		if err == storage.ErrUrlExists {
			http.Error(w, "url already exists", http.StatusBadRequest)
			h.log.Error(err)
			return
		}

		http.Error(w, "error with saving url"+fmt.Sprintf("%s", err), http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	resp, err := json.Marshal(map[string]int64{
		"id": id,
	})
	if err != nil {
		http.Error(w, "error with creating response", http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	h.log.Info(fmt.Sprintf("add new [url]: (%s), [alias]: (%s)", url.URL, url.Alias))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *Handler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alias := vars["alias"]

	err := h.urlService.Delete(context.Background(), alias)
	if err != nil {
		if err == storage.ErrURLNotFound {
			http.Error(w, "url not found", http.StatusNotFound)
			h.log.Error(err)
			return
		}

		http.Error(w, "error with deleting url", http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	h.log.Info(fmt.Sprintf("[alias]: (%s) - was deleted", alias))
	w.WriteHeader(http.StatusNoContent)

}
