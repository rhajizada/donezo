package itemsbytag_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/itemsbytag"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

func newItemsByTagMenu(t *testing.T) (itemsbytag.MenuModel, *service.Service, func()) {
	t.Helper()
	svc, cleanup := testutil.NewTestService(t)
	ctx := testutil.MustContext()

	board, err := svc.CreateBoard(ctx, "Inbox")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}
	item, err := svc.CreateItem(ctx, board, "task", "desc")
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	item.Tags = []string{"work"}
	item, err = svc.UpdateItem(ctx, item)
	if err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}

	tagCount, err := svc.CountItemsByTag(ctx, "work")
	if err != nil {
		t.Fatalf("CountItemsByTag: %v", err)
	}
	parent := tags.NewModel(ctx, svc)
	parent.List.SetItems(tags.NewList([]tags.Item{tags.NewItem("work", tagCount)}))
	parent.List.Select(0)

	menu := itemsbytag.New(ctx, svc, &parent)
	menu.List.SetItems(itemsbytag.NewList(&[]service.Item{*item}))
	menu.List.Select(0)
	return menu, svc, cleanup
}

func TestListItemsWithoutSelectedTagReturnsErrorMsg(t *testing.T) {
	menu, _, cleanup := newItemsByTagMenu(t)
	defer cleanup()

	menu.Parent.List.SetItems(tags.NewList(nil))

	msg := menu.ListItems()()
	if _, ok := msg.(itemsbytag.ErrorMsg); !ok {
		t.Fatalf("expected ErrorMsg, got %T", msg)
	}
}

func TestRenameUpdateTagsAndToggleCommands(t *testing.T) {
	menu, _, cleanup := newItemsByTagMenu(t)
	defer cleanup()

	menu.Context.Title = "renamed"
	menu.Context.Desc = "renamed desc"

	msg := menu.RenameItem()()
	renameMsg, ok := msg.(itemsbytag.RenameItemMsg)
	if !ok {
		t.Fatalf("expected RenameItemMsg, got %T", msg)
	}
	if renameMsg.Error != nil {
		t.Fatalf("RenameItem error: %v", renameMsg.Error)
	}

	menu.Context.Title = "one, two"
	msg = menu.UpdateTags()()
	updateTagsMsg, ok := msg.(itemsbytag.UpdateTagsMsg)
	if !ok {
		t.Fatalf("expected UpdateTagsMsg, got %T", msg)
	}
	if updateTagsMsg.Error != nil {
		t.Fatalf("UpdateTags error: %v", updateTagsMsg.Error)
	}
	if len(updateTagsMsg.Item.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %+v", updateTagsMsg.Item.Tags)
	}

	msg = menu.ToggleComplete()()
	toggleMsg, ok := msg.(itemsbytag.ToggleItemMsg)
	if !ok {
		t.Fatalf("expected ToggleItemMsg, got %T", msg)
	}
	if toggleMsg.Error != nil {
		t.Fatalf("ToggleComplete error: %v", toggleMsg.Error)
	}
	if !toggleMsg.Item.Completed {
		t.Fatalf("expected item marked complete")
	}
}

func TestUpdateTagsValidationError(t *testing.T) {
	menu, _, cleanup := newItemsByTagMenu(t)
	defer cleanup()

	menu.Context.Title = "ok, "
	msg := menu.UpdateTags()()
	updateMsg, ok := msg.(itemsbytag.UpdateTagsMsg)
	if !ok {
		t.Fatalf("expected UpdateTagsMsg, got %T", msg)
	}
	if updateMsg.Error == nil {
		t.Fatalf("expected validation error")
	}
}

func TestHandleInputStateTransitions(t *testing.T) {
	menu, _, cleanup := newItemsByTagMenu(t)
	defer cleanup()

	menu.Context.State = itemsbytag.RenameItemNameState
	menu.Input.SetValue("new title")
	_, _ = menu.HandleInputState(tea.KeyPressMsg{Code: tea.KeyEnter})
	if menu.Context.State != itemsbytag.RenameItemDescState {
		t.Fatalf("expected rename desc state, got %v", menu.Context.State)
	}

	menu.Input.SetValue("new description")
	_, cmds := menu.HandleInputState(tea.KeyPressMsg{Code: tea.KeyEnter})
	if menu.Context.State != itemsbytag.DefaultState {
		t.Fatalf("expected default state, got %v", menu.Context.State)
	}
	if len(cmds) == 0 {
		t.Fatalf("expected commands from rename completion")
	}

	menu.Context.State = itemsbytag.UpdateTagsState
	menu.Input.SetValue("a, b")
	_, cmds = menu.HandleInputState(tea.KeyPressMsg{Code: tea.KeyEnter})
	if menu.Context.State != itemsbytag.DefaultState {
		t.Fatalf("expected default state after update tags, got %v", menu.Context.State)
	}
	if len(cmds) == 0 {
		t.Fatalf("expected commands from update tags completion")
	}

	menu.Context.State = itemsbytag.RenameItemNameState
	_, _ = menu.HandleInputState(tea.KeyPressMsg{Code: tea.KeyEsc})
	if menu.Context.State != itemsbytag.DefaultState {
		t.Fatalf("expected escape to return default state, got %v", menu.Context.State)
	}
}

func TestUpdateWithListItemsMsgReplacesList(t *testing.T) {
	menu, _, cleanup := newItemsByTagMenu(t)
	defer cleanup()

	newItems := &[]service.Item{{Item: service.Item{}.Item, Tags: []string{"x"}}}
	(*newItems)[0].Title = "replacement"
	(*newItems)[0].Description = "desc"

	model, _ := menu.Update(itemsbytag.ListItemsMsg{Items: newItems})
	updated := model.(itemsbytag.MenuModel)
	if len(updated.List.Items()) != 1 {
		t.Fatalf("expected list replaced with 1 item, got %d", len(updated.List.Items()))
	}
}
