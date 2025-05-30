package itemsbyboard

import "github.com/rhajizada/donezo/internal/service"

type ErrorMsg struct {
	Error error
}

type ListItemsMsg struct {
	Items *[]service.Item
}

type CreateItemMsg struct {
	Item  *service.Item
	Error error
}

type RenameItemMsg struct {
	Item  *service.Item
	Error error
}

type UpdateTagsMsg struct {
	Item  *service.Item
	Error error
}

type ToggleItemMsg struct {
	Item  *service.Item
	Error error
}

type DeleteItemMsg struct {
	Item  *service.Item
	Error error
}
