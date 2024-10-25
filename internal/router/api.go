package router

import (
	"net/http"

	_ "github.com/rhajizada/donezo/docs"
	"github.com/rhajizada/donezo/internal/handler"
)

func RegisterApiRoutes(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /boards", h.ListBoards)
	router.HandleFunc("GET /boards/{boardId}", h.GetBoardByID)
	router.HandleFunc("POST /boards", h.CreateBoard)
	router.HandleFunc("DELETE /boards/{boardId}", h.DeleteBoardByID)
	router.HandleFunc("PUT /boards/{boardId}", h.UpdateBoardByID)
	router.HandleFunc("GET /boards/{boardId}/items", h.ListItemsByBoardID)
	router.HandleFunc("POST /boards/{boardId}/items", h.AddItemToBoardByID)
	router.HandleFunc("PUT /boards/{boardId}/items/{itemId}", h.UpdateItemById)
	router.HandleFunc("DELETE /boards/{boardId}/items/{itemId}", h.DeleteItemByID)
	return router
}
