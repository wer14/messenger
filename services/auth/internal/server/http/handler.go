package http

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type ReadyChecker interface {
	IsReady() bool
}

type Handler struct {
	checker ReadyChecker
}

func NewHandler(checker ReadyChecker) *Handler {
	return &Handler{
		checker: checker,
	}
}

func (h *Handler) Health(c echo.Context) error {
	log.Println("[probe] /health check received")

	return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}

func (h *Handler) Ready(c echo.Context) error {
	log.Println("[probe] /ready check received")

	if !h.checker.IsReady() {
		log.Println("[probe] /ready - NOT_SERVING")

		return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unavalaible"})
	}

	log.Println("[probe] /ready - SERVING")

	return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}
