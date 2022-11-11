package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/tuanp/go-gin-boilerplate/services/app1/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init(e *echo.Echo) {
	//g := e.Group("/v1")
}
