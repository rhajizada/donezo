package navigation

// View represents the active page in the TUI.
type View int

const (
	ViewBoards View = iota
	ViewTags
	ViewItemsByBoard
	ViewItemsByTag
)

// SwitchMainViewMsg requests swapping between the root menus (boards <-> tags).
type SwitchMainViewMsg struct {
	View View
}

// OpenBoardItemsMsg requests opening the items view for the selected board.
type OpenBoardItemsMsg struct{}

// OpenTagItemsMsg requests opening the items view for the selected tag.
type OpenTagItemsMsg struct{}

// BackMsg requests returning to the previous view (from detail to its parent menu).
type BackMsg struct{}

// BoardDeltaMsg requests moving the board selection by Delta and refreshing dependent views.
type BoardDeltaMsg struct {
	Delta int
}

// TagDeltaMsg requests moving the tag selection by Delta and refreshing dependent views.
type TagDeltaMsg struct {
	Delta int
}
