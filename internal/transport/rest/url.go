package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/f0xg0sasha/url_short/internal/domain"
	"github.com/f0xg0sasha/url_short/internal/storage"
	"github.com/gorilla/mux"
)

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alias := vars["alias"]

	OriginalURL, err := h.urlService.GetURL(alias)
	if err != nil {
		if err == storage.ErrURLNotFound {
			http.Error(w, "error", http.StatusNotFound)
			return
		}

		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, OriginalURL, http.StatusMovedPermanently)
}

func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	url := &domain.RequestURL{}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBytes, url)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	id, err := h.urlService.SaveURL(url.URL, url.Alias)
	if err != nil {
		if err == storage.ErrUrlExists {
			http.Error(w, "url already exists", http.StatusBadRequest)
			return
		}

		http.Error(w, "error with saving url"+fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(map[string]int64{
		"id": id,
	})
	if err != nil {
		http.Error(w, "error with creating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *Handler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alias := vars["alias"]

	err := h.urlService.DeleteURL(alias)
	if err != nil {
		if err == storage.ErrURLNotFound {
			http.Error(w, "url not found", http.StatusNotFound)
			log.Print(err)
			return
		}

		http.Error(w, "error with deleting url", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}
