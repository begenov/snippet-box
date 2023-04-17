package controller

import (
	"net/http"

	"github.com/begenov/snippet-box/internal/service"
)

type AuthHandlers interface {
	Router() http.Handler
}

type Handler struct {
	services service.AuthApplication
}

func NewHandler(services service.AuthApplication) AuthHandlers {
	return &Handler{
		services: services,
	}
}
