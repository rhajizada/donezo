package boards

import "github.com/rhajizada/donezo/internal/service"

type ErrorMsg struct {
	Error error
}

type ListBoardsMsg struct {
	Boards *[]service.Board
}

type CreateBoardMsg struct {
	Board *service.Board
	Error error
}

type RenameBoardMsg struct {
	Board *service.Board
	Error error
}

type DeleteBoardMsg struct {
	Board *service.Board
	Error error
}
