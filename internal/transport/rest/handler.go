package rest

import "github.com/gorilla/mux"

type URLService interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	GetURL(alias string) (string, error)
	DeleteURL(alias string) error
}

type Handler struct {
	urlService URLService
}

func NewHandler(u URLService) *Handler {
	return &Handler{
		urlService: u,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()

	url := r.PathPrefix("/url").Subrouter()
	{
		url.HandleFunc("/{alias:[a-z]+}", h.GetURL).Methods("GET")
		url.HandleFunc("/", h.CreateURL).Methods("POST")
		url.HandleFunc("/{alias:[a-z]+}", h.DeleteURL).Methods("DELETE")
	}

	return r
}
