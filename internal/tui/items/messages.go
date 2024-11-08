package items

import "github.com/rhajizada/donezo/pkg/client"

type ErrorMsg struct {
	Error error
}

type ListItemsMsg struct {
	Items *[]client.Item
}

type CreateItemMsg struct {
	Item  *client.Item
	Error error
}

type RenameItemMsg struct {
	Item  *client.Item
	Error error
}

type ToggleItemMsg struct {
	Item  *client.Item
	Error error
}

type DeleteItemMsg struct {
	Item  *client.Item
	Error error
}
