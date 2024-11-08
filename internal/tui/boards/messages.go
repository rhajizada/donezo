package boards

import "github.com/rhajizada/donezo/pkg/client"

type ErrorMsg struct {
	Error error
}

type ListBoardsMsg struct {
	Boards *[]client.Board
}

type CreateBoardMsg struct {
	Board *client.Board
	Error error
}

type RenameBoardMsg struct {
	Board *client.Board
	Error error
}

type DeleteBoardMsg struct {
	Board *client.Board
	Error error
}
