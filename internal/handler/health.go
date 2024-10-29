package handler

import (
	"encoding/json"
	"net/http"
)

// Healthz godoc
// @Summary Check health
// @Description Get service health status
// @Produce json
// @Success 200 {array} StatusResponse
// @Failure 500 {object} string
// @Router /healthz [get]
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(DefaultHealthResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
