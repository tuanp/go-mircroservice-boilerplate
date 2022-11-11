package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/tuanp/go-gin-boilerplate/config"
	v1 "github.com/tuanp/go-gin-boilerplate/services/app1/internal/handlers/v1"
	"github.com/tuanp/go-gin-boilerplate/services/app1/internal/services"
	"net/http"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init(cfg *config.Config) *echo.Echo {
	e := echo.New()

	// Init router
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	h.initAPI(e)

	return e
}

func (h *Handler) initAPI(e *echo.Echo) {
	handlerV1 := v1.NewHandler(h.services)
	handlerV1.Init(e)
}
