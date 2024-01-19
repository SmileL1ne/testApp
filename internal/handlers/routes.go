package handlers

import (
	"log/slog"
	"net/http"
	"testApp/internal/service"

	"github.com/julienschmidt/httprouter"
)

type routes struct {
	logger  *slog.Logger
	service *service.Service
}

func NewRouter(logger *slog.Logger, s *service.Service) http.Handler {
	router := httprouter.New()

	r := &routes{
		logger:  logger,
		service: s,
	}

	router.HandlerFunc(http.MethodGet, "/users", r.getAllUsers)
	router.HandlerFunc(http.MethodDelete, "/users/:id", r.deleteUser)
	router.HandlerFunc(http.MethodPut, "/users/:id", r.updateUser)
	router.HandlerFunc(http.MethodPost, "/users", r.saveUser)

	return router
}
