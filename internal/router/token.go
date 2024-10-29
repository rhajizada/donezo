package router

import (
	"net/http"

	_ "github.com/rhajizada/donezo/docs"
	"github.com/rhajizada/donezo/internal/handler"
)

func RegisterTokenRoutes(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /token/refresh", h.RefreshToken)
	router.HandleFunc("GET /token/validate", h.ValidateToken)
	return router
}
