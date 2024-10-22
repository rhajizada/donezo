package router

import (
	"net/http"

	"github.com/rhajizada/donezo/internal/handler"
)

func RegisterApiRoutes(h *handler.Handler) *http.ServeMux {
	boardsRouter := http.NewServeMux()
	boardsRouter.HandleFunc("GET /boards", h.ListBoards)
	boardsRouter.HandleFunc("GET /boards/{boardId}", h.GetBoardByID)
	boardsRouter.HandleFunc("POST /boards", h.CreateBoard)
	boardsRouter.HandleFunc("DELETE /boards/{boardId}", h.DeleteBoardByID)
	boardsRouter.HandleFunc("PUT /boards/{boardId}", h.UpdateBoardByID)
	boardsRouter.HandleFunc("GET /boards/{boardId}/items", h.ListItemsByBoardID)
	boardsRouter.HandleFunc("POST /boards/{boardId}/items", h.AddItemToBoardByID)
	boardsRouter.HandleFunc("PUT /boards/{boardId}/items/{itemId}", h.UpdateItemById)
	boardsRouter.HandleFunc("DELETE /boards/{boardId}/items/{itemId}", h.DeleteItemByID)
	router := http.NewServeMux()
	router.Handle("/api/", http.StripPrefix("/api", boardsRouter))
	return router
}
