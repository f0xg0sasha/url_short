package rest

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type URLService interface {
	Fetch(ctx context.Context, alias string) (string, error)
	Create(ctx context.Context, url string, alias string) (int64, error)
	Delete(ctx context.Context, alias string) error
}

type Handler struct {
	log        *logrus.Logger
	urlService URLService
}

func NewHandler(log *logrus.Logger, urlService URLService) *Handler {
	return &Handler{
		log:        log,
		urlService: urlService,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()

	url := r.PathPrefix("/url").Subrouter()
	{
		url.HandleFunc("/{alias:[a-z0-9]+}", h.GetURL).Methods("GET")
		url.HandleFunc("/", h.CreateURL).Methods("POST")
		url.HandleFunc("/{alias:[a-z0-9]+}", h.DeleteURL).Methods("DELETE")
	}

	return r
}
