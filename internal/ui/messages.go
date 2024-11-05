package ui

import "github.com/rhajizada/donezo/pkg/client"

// Message types
type errMsg struct {
	err error
}

type boardsLoadedMsg struct{}

type itemsLoadedMsg struct {
	items []client.Item
}

type updateItemMsg struct {
	err error
}

type createItemMsg struct {
	item client.Item
	err  error
}
