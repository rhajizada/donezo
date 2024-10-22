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
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	var input BoardRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		msg := fmt.Sprintf("error decoding JSON: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	data, err := h.Repo.CreateBoard(r.Context(), input.Name)
	if err != nil {
		msg := fmt.Sprintf("failed creating board %v: %v", input.Name, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetBoardByID(w http.ResponseWriter, r *http.Request) {
	boardId := r.PathValue("boardId")
	if boardId == "" {
		http.Error(w, "missing 'id' parameter", http.StatusBadRequest)
		return
	}
	boardIdAsInt, err := strconv.ParseInt(boardId, 10, 64)
	if err != nil {
		http.Error(w, "cannot parse 'id' parameter into integer", http.StatusInternalServerError)
		return
	}
	data, err := h.Repo.GetBoardByID(r.Context(), boardIdAsInt)
	if err != nil {
		msg := fmt.Sprintf("failed fetching board %d: %v", boardIdAsInt, err)
		http.Error(w, msg, http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateBoardByID(w http.ResponseWriter, r *http.Request) {
	boardId := r.PathValue("boardId")
	var input BoardRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if boardId == "" {
		http.Error(w, "missing 'id' parameter", http.StatusBadRequest)
		return
	}
	boardIdAsInt, err := strconv.ParseInt(boardId, 10, 64)
	if err != nil {
		http.Error(w, "cannot parse 'id' parameter into integer", http.StatusInternalServerError)
		return
	}
	b := repository.UpdateBoardByIDParams{
		Name: input.Name,
		ID:   boardIdAsInt,
	}
	data, err := h.Repo.UpdateBoardByID(r.Context(), b)
	if err != nil {
		msg := fmt.Sprintf("failed updating board %d: %v", boardIdAsInt, err)
		http.Error(w, msg, http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteBoardByID(w http.ResponseWriter, r *http.Request) {
	boardId := r.PathValue("boardId")
	var input BoardRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if boardId == "" {
		http.Error(w, "missing 'id' parameter", http.StatusBadRequest)
		return
	}
	boardIdAsInt, err := strconv.ParseInt(boardId, 10, 64)
	if err != nil {
		http.Error(w, "cannot parse 'id' parameter into integer", http.StatusInternalServerError)
		return
	}

	_, err = h.Repo.GetBoardByID(r.Context(), boardIdAsInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.Repo.DeleteBoardByID(r.Context(), boardIdAsInt)
	msg := fmt.Sprintf("succesfully deleted board %v", boardId)
	w.Write([]byte(msg))
}

func (h *Handler) ListItemsByBoardID(w http.ResponseWriter, r *http.Request) {
	boardId := r.PathValue("boardId")

	if boardId == "" {
		http.Error(w, "missing 'id' parameter", http.StatusBadRequest)
		return
	}
	boardIdAsInt, err := strconv.ParseInt(boardId, 10, 64)
	if err != nil {
		http.Error(w, "cannot parse 'id' parameter into integer", http.StatusInternalServerError)
		return
	}

	data, err := h.Repo.ListItemsByBoardID(r.Context(), boardIdAsInt)
	if err != nil {
		msg := fmt.Sprintf("failed fetching items for board %d: %v", boardIdAsInt, err)
		http.Error(w, msg, http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddItemToBoardByID(w http.ResponseWriter, r *http.Request) {
	boardId := r.PathValue("boardId")

	if boardId == "" {
		http.Error(w, "missing 'id' parameter", http.StatusBadRequest)
		return
	}
	boardIdAsInt, err := strconv.ParseInt(boardId, 10, 64)
	if err != nil {
		http.Error(w, "cannot parse 'id' parameter into integer", http.StatusInternalServerError)
		return
	}

	var input CreateItemRequest
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		msg := fmt.Sprintf("error decoding JSON: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	i := repository.CreateItemParams{
		BoardID:     boardIdAsInt,
		Title:       input.Title,
		Description: input.Description,
	}

	data, err := h.Repo.CreateItem(r.Context(), i)
	if err != nil {
		msg := fmt.Sprintf("failed adding item %v to board %d: %v", input.Title, boardIdAsInt, err)
		http.Error(w, msg, http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateItemById(w http.ResponseWriter, r *http.Request) {
	boardId := r.PathValue("boardId")

	if boardId == "" {
		http.Error(w, "missing 'boardId' parameter", http.StatusBadRequest)
		return
	}
	boardIdAsInt, err := strconv.ParseInt(boardId, 10, 64)
	if err != nil {
		http.Error(w, "cannot parse 'boardId' parameter into integer", http.StatusInternalServerError)
		return
	}
	_, err = h.Repo.GetBoardByID(r.Context(), boardIdAsInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	itemId := r.PathValue("itemId")
	if itemId == "" {
		http.Error(w, "missing 'itemId' parameter", http.StatusBadRequest)
		return
	}
	itemIdAsInt, err := strconv.ParseInt(itemId, 10, 64)
	if err != nil {
		http.Error(w, "cannot parse 'itemId' parameter into integer", http.StatusInternalServerError)
		return
	}

	var input UpdateItemRequest
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		msg := fmt.Sprintf("error decoding JSON: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	i := repository.UpdateItemByIDParams{
		Title:       input.Title,
		Description: input.Description,
		Completed:   input.Completed,
		ID:          itemIdAsInt,
	}

	data, err := h.Repo.UpdateItemByID(r.Context(), i)
	if err != nil {
		msg := fmt.Sprintf("failed updating item %d: %v", itemIdAsInt, err)
		http.Error(w, msg, http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteItemByID(w http.ResponseWriter, r *http.Request) {
	boardId := r.PathValue("boardId")

	if boardId == "" {
		http.Error(w, "missing 'boardId' parameter", http.StatusBadRequest)
		return
	}
	boardIdAsInt, err := strconv.ParseInt(boardId, 10, 64)
	if err != nil {
		http.Error(w, "cannot parse 'boardId' parameter into integer", http.StatusInternalServerError)
		return
	}
	_, err = h.Repo.GetBoardByID(r.Context(), boardIdAsInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	itemId := r.PathValue("itemId")
	if itemId == "" {
		http.Error(w, "missing 'itemId' parameter", http.StatusBadRequest)
		return
	}
	itemIdAsInt, err := strconv.ParseInt(itemId, 10, 64)
	if err != nil {
		http.Error(w, "cannot parse 'itemId' parameter into integer", http.StatusInternalServerError)
		return
	}

	_, err = h.Repo.GetItemByID(r.Context(), itemIdAsInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.Repo.DeleteItemByID(r.Context(), itemIdAsInt)
	msg := fmt.Sprintf("succesfully deleted item %v", itemId)
	w.Write([]byte(msg))
}
