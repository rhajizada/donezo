package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rhajizada/donezo/internal/repository"
)

type Handler struct {
	Repo repository.Queries
}

func New(r *repository.Queries) *Handler {
	return &Handler{
		Repo: *r,
	}
}

func (h *Handler) ListBoards(w http.ResponseWriter, r *http.Request) {
	data, err := h.Repo.ListBoards(r.Context())
	if err != nil {
		msg := fmt.Sprintf("failed fetching boards : %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	var input CreateBoardInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		msg := fmt.Sprintf("error decoding JSON: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
	}

	data, err := h.Repo.CreateBoard(r.Context(), input.Name)
	if err != nil {
		msg := fmt.Sprintf("failed creating board %v", input.Name)
		http.Error(w, msg, http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetBoardByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing 'id' parameter", http.StatusBadRequest)
		return
	}
	idAsInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "cannot parse 'id' parameter into integer", http.StatusInternalServerError)
		return
	}
	data, err := h.Repo.GetBoardByID(r.Context(), idAsInt)
	if err != nil {
		msg := fmt.Sprintf("failed fetching board %d", idAsInt)
		http.Error(w, msg, http.StatusNotFound)
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
