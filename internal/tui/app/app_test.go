package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/itemsbytag"
	"github.com/rhajizada/donezo/internal/tui/navigation"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

func seedBoard(t *testing.T, svc *service.Service, name string) *service.Board {
	t.Helper()
	board, err := svc.CreateBoard(testutil.MustContext(), name)
	require.NoError(t, err)
	return board
}

func TestAppNavigatesBetweenBoardsAndItems(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "open board items then navigate back"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			board := seedBoard(t, svc, "Inbox")
			_, err := svc.CreateItem(testutil.MustContext(), board, "task", "desc")
			require.NoError(t, err)

			ctx := testutil.MustContext()
			m := New(ctx, svc)
			m.boards.List.SetItems(boards.NewList(&[]service.Board{*board}))

			model, _ := m.Update(navigation.OpenBoardItemsMsg{})
			am, ok := model.(AppModel)
			require.True(t, ok)
			assert.Equal(t, navigation.ViewItemsByBoard, am.active)
			assert.NotNil(t, am.itemsByBoard)

			model, _ = am.Update(navigation.BackMsg{})
			am, ok = model.(AppModel)
			require.True(t, ok)
			assert.Equal(t, navigation.ViewBoards, am.active)
		})
	}
}

func TestBoardsEscDoesNotQuit(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "esc keeps boards view active"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			board := seedBoard(t, svc, "Inbox")
			ctx := testutil.MustContext()
			menu := boards.New(ctx, svc)
			menu.List.SetItems(boards.NewList(&[]service.Board{*board}))

			model, cmd := menu.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
			if cmd != nil {
				_, ok := cmd().(tea.QuitMsg)
				assert.False(t, ok)
			}

			updated, ok := model.(boards.MenuModel)
			require.True(t, ok)
			assert.Equal(t, boards.DefaultState, updated.State)
		})
	}
}

func TestTagsEscDoesNotQuit(t *testing.T) {
	tests := []struct {
		name string
		tag  string
	}{
		{name: "esc keeps tags view active", tag: "inbox"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			board := seedBoard(t, svc, "Inbox")
			item, err := svc.CreateItem(testutil.MustContext(), board, "task", "desc")
			require.NoError(t, err)
			item.Tags = []string{tt.tag}
			_, err = svc.UpdateItem(testutil.MustContext(), item)
			require.NoError(t, err)

			tagCount, err := svc.CountItemsByTag(testutil.MustContext(), tt.tag)
			require.NoError(t, err)

			ctx := testutil.MustContext()
			menu := tags.NewModel(ctx, svc)
			menu.List.SetItems(tags.NewList([]tags.Item{tags.NewItem(tt.tag, tagCount)}))

			model, cmd := menu.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
			if cmd != nil {
				_, ok := cmd().(tea.QuitMsg)
				assert.False(t, ok)
			}

			_, ok := model.(tags.MenuModel)
			assert.True(t, ok)
		})
	}
}

func TestBoardsTabSwitchesToTags(t *testing.T) {
	tests := []struct {
		name string
		tag  string
	}{
		{name: "tab from boards switches to tags", tag: "work"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			board := seedBoard(t, svc, "Inbox")
			item, err := svc.CreateItem(ctx, board, "task", "desc")
			require.NoError(t, err)
			item.Tags = []string{tt.tag}
			_, err = svc.UpdateItem(ctx, item)
			require.NoError(t, err)

			m := New(ctx, svc)
			m.boards.List.SetItems(boards.NewList(&[]service.Board{*board}))

			model, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
			appModel := model.(AppModel)
			require.NotNil(t, cmd)

			switchMsg := cmd()
			model, cmd = appModel.Update(switchMsg)
			appModel = model.(AppModel)
			if cmd != nil {
				if msg := cmd(); msg != nil {
					model, _ = appModel.Update(msg)
					appModel = model.(AppModel)
				}
			}

			assert.Equal(t, navigation.ViewTags, appModel.active)
			assert.Len(t, appModel.tags.List.Items(), 1)
		})
	}
}

func TestItemsByBoardTabAndShiftTabCycleBoards(t *testing.T) {
	tests := []struct {
		name          string
		delta         int
		wantIndex     int
		wantBoardName string
	}{
		{name: "tab selects next board", delta: 1, wantIndex: 1, wantBoardName: "Two"},
		{name: "shift tab selects previous board", delta: -1, wantIndex: 0, wantBoardName: "One"},
	}

	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()
	ctx := testutil.MustContext()
	board1 := seedBoard(t, svc, "One")
	board2 := seedBoard(t, svc, "Two")
	_, err := svc.CreateItem(ctx, board1, "a", "")
	require.NoError(t, err)
	_, err = svc.CreateItem(ctx, board2, "b", "")
	require.NoError(t, err)

	m := New(ctx, svc)
	m.boards.List.SetItems(boards.NewList(&[]service.Board{*board1, *board2}))
	model, _ := m.Update(navigation.OpenBoardItemsMsg{})
	appModel := model.(AppModel)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextModel, cmd := appModel.Update(navigation.BoardDeltaMsg{Delta: tt.delta})
			updated := nextModel.(AppModel)
			if cmd != nil {
				if msg := cmd(); msg != nil {
					nextModel, _ = updated.Update(msg)
					updated = nextModel.(AppModel)
				}
			}
			assert.Equal(t, tt.wantIndex, updated.boards.List.Index())
			assert.Equal(t, tt.wantBoardName, updated.itemsByBoard.List.Title)
			appModel = updated
		})
	}
}

func TestItemsByTagNextAndPrevious(t *testing.T) {
	tests := []struct {
		name         string
		delta        int
		wantIndex    int
		wantTagTitle string
	}{
		{name: "next tag moves forward", delta: 1, wantIndex: 1, wantTagTitle: "beta"},
		{name: "previous tag moves backward", delta: -1, wantIndex: 0, wantTagTitle: "alpha"},
	}

	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()
	ctx := testutil.MustContext()
	board := seedBoard(t, svc, "Inbox")
	itemA, err := svc.CreateItem(ctx, board, "a", "")
	require.NoError(t, err)
	itemA.Tags = []string{"alpha"}
	_, err = svc.UpdateItem(ctx, itemA)
	require.NoError(t, err)
	itemB, err := svc.CreateItem(ctx, board, "b", "")
	require.NoError(t, err)
	itemB.Tags = []string{"beta"}
	_, err = svc.UpdateItem(ctx, itemB)
	require.NoError(t, err)

	m := New(ctx, svc)
	tagCountAlpha, err := svc.CountItemsByTag(ctx, "alpha")
	require.NoError(t, err)
	tagCountBeta, err := svc.CountItemsByTag(ctx, "beta")
	require.NoError(t, err)
	m.tags.List.SetItems(tags.NewList([]tags.Item{
		tags.NewItem("alpha", tagCountAlpha),
		tags.NewItem("beta", tagCountBeta),
	}))

	model, _ := m.Update(navigation.OpenTagItemsMsg{})
	appModel := model.(AppModel)
	appModel.tags.List.Select(0)
	itemTagModel := itemsbytag.New(ctx, svc, appModel.tags)
	appModel.itemsByTag = &itemTagModel
	appModel.active = navigation.ViewItemsByTag

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextModel, cmd := appModel.Update(navigation.TagDeltaMsg{Delta: tt.delta})
			updated := nextModel.(AppModel)
			if cmd != nil {
				if msg := cmd(); msg != nil {
					nextModel, _ = updated.Update(msg)
					updated = nextModel.(AppModel)
				}
			}
			assert.Equal(t, tt.wantIndex, updated.tags.List.Index())
			assert.Equal(t, tt.wantTagTitle, updated.itemsByTag.List.Title)
			appModel = updated
		})
	}
}
