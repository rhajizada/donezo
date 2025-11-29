package tags

import (
	"testing"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
)

func TestCopyAndDeleteTag(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

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
	if _, err := svc.UpdateItem(ctx, item); err != nil {
		t.Fatalf("UpdateItem set tag: %v", err)
	}

	menu := NewModel(ctx, svc)
	msg := menu.ListTags()()
	model, _ := menu.Update(msg)
	menu = model.(MenuModel)
	menu.List.Select(0)

	itemsForTag, err := svc.ListItemsByTag(ctx, "work")
	if err != nil {
		t.Fatalf("ListItemsByTag: %v", err)
	}
	expected := service.ItemsToMarkdown("work", *itemsForTag)

	var captured []byte
	prevWrite := writeClipboardText
	writeClipboardText = func(data []byte) { captured = append([]byte{}, data...) }
	defer func() { writeClipboardText = prevWrite }()

	if cmd := menu.Copy(); cmd != nil {
		cmd()
	}
	if string(captured) != expected {
		t.Fatalf("unexpected clipboard markdown\nwant:\n%s\n\ngot:\n%s", expected, string(captured))
	}

	// Delete tag and ensure list refreshes.
	delCmd := menu.DeleteTag()
	msg = delCmd()
	model, _ = menu.Update(msg)
	menu = model.(MenuModel)

	refresh := menu.ListTags()
	msg = refresh()
	model, _ = menu.Update(msg)
	menu = model.(MenuModel)

	if len(menu.List.Items()) != 0 {
		t.Fatalf("expected no tags after delete, got %d", len(menu.List.Items()))
	}
	count, err := svc.CountItemsByTag(ctx, "work")
	if err != nil {
		t.Fatalf("CountItemsByTag: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected tag removal in service, count=%d", count)
	}
}
