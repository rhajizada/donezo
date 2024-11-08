package items

type InputState uint8

const (
	DefaultState InputState = iota
	CreateItemNameState
	CreateItemDescState
	RenameItemNameState
	RenameItemDescState
)

type InputContext struct {
	State InputState
	Title string
	Desc  string
}

func NewInputContext() *InputContext {
	return &InputContext{
		State: DefaultState,
		Title: "",
		Desc:  "",
	}
}
