package boards

type InputState uint8

const (
	DefaultState InputState = iota
	CreateBoardState
	RenameBoardState
)
