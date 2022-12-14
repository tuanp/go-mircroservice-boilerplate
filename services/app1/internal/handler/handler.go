package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tuanp/go-mircroservice-boilerplate/pkg/config"
	"github.com/tuanp/go-mircroservice-boilerplate/pkg/logger"
	v1 "github.com/tuanp/go-mircroservice-boilerplate/services/app1/internal/handler/v1"
	"github.com/tuanp/go-mircroservice-boilerplate/services/app1/internal/service"
)

type Handler struct {
	services *service.Services
	logger   logger.Logger
}

func NewHandler(services *service.Services, logger logger.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
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
