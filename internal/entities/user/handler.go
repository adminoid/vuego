package user

import (
	"context"
	"encoding/json"
	"github.com/adminoid/vuego/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type handler struct {
	logger     *logging.Logger
	repository RepositoryUser
}

type Handler interface {
	Register(router *httprouter.Router)
}

func NewHandler(repository RepositoryUser, logger *logging.Logger) Handler {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET("/users", h.GetList)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	all, err := h.repository.FindAll(context.TODO())
	if err != nil {
		http.Error(w, "Forbidden", 400)
	}

	allBytes, err := json.Marshal(all)
	if err != nil {
		http.Error(w, "Parsing json error", 500)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)
}
