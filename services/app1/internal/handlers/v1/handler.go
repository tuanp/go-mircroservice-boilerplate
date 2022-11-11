package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/tuanp/go-gin-boilerplate/services/app1/internal/services"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init(e *echo.Echo) {
	//g := e.Group("/v1")
}
