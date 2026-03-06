package itemlist_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/rhajizada/donezo/internal/tui/itemlist"
)

type extItem struct {
	title  string
	hidden bool
}

func (i extItem) FilterValue() string { return i.title }
func (i extItem) HideValue() bool     { return i.hidden }

type extDelegate struct {
	height  int
	spacing int
}

func (d extDelegate) Render(w io.Writer, _ itemlist.Model, _ int, item itemlist.Item) {
	fmt.Fprint(w, item.FilterValue())
}

func (d extDelegate) Height() int {
	if d.height == 0 {
		return 1
	}
	return d.height
}
func (d extDelegate) Spacing() int { return d.spacing }
func (d extDelegate) Update(tea.Msg, *itemlist.Model) tea.Cmd {
	return nil
}

func newModel(items []itemlist.Item) itemlist.Model {
	return itemlist.New(items, extDelegate{height: 1, spacing: 0}, 80, 20)
}

func TestPublicSettersAndGetters(t *testing.T) {
	m := newModel([]itemlist.Item{extItem{title: "a"}})

	m.SetShowTitle(false)
	m.SetShowFilter(false)
	m.SetShowStatusBar(false)
	m.SetShowPagination(false)
	m.SetShowHelp(false)
	m.SetFilteringEnabled(true)
	m.SetStatusBarItemName("task", "tasks")
	m.SetSize(70, 15)

	if m.ShowTitle() {
		t.Fatalf("expected title hidden")
	}
	if m.ShowFilter() {
		t.Fatalf("expected filter hidden")
	}
	if m.ShowStatusBar() {
		t.Fatalf("expected status bar hidden")
	}
	if m.ShowPagination() {
		t.Fatalf("expected pagination hidden")
	}
	if m.ShowHelp() {
		t.Fatalf("expected help hidden")
	}
	if !m.FilteringEnabled() {
		t.Fatalf("expected filtering enabled")
	}
	singular, plural := m.StatusBarItemName()
	if singular != "task" || plural != "tasks" {
		t.Fatalf("unexpected status names: %q, %q", singular, plural)
	}
	if m.Width() != 70 || m.Height() != 15 {
		t.Fatalf("unexpected size %dx%d", m.Width(), m.Height())
	}
}

func TestItemMutationAndSelection(t *testing.T) {
	m := newModel([]itemlist.Item{extItem{title: "a"}, extItem{title: "b"}})

	m.InsertItem(1, extItem{title: "inserted"})
	if len(m.Items()) != 3 {
		t.Fatalf("expected 3 items, got %d", len(m.Items()))
	}

	m.SetItem(0, extItem{title: "first"})
	m.Select(1)
	if m.Index() != 1 || m.GlobalIndex() != 1 {
		t.Fatalf("unexpected selected index local=%d global=%d", m.Index(), m.GlobalIndex())
	}

	m.RemoveItem(1)
	if len(m.Items()) != 2 {
		t.Fatalf("expected 2 items after remove, got %d", len(m.Items()))
	}

	m.ResetSelected()
	selected := m.SelectedItem()
	if selected == nil {
		t.Fatalf("expected selected item after reset")
	}
}

func TestCursorPagingAndInfiniteScrolling(t *testing.T) {
	items := []itemlist.Item{
		extItem{title: "1"},
		extItem{title: "2"},
		extItem{title: "3"},
	}
	m := itemlist.New(items, extDelegate{height: 1, spacing: 0}, 40, 2)
	m.InfiniteScrolling = true

	m.CursorUp()
	if m.Index() != len(items)-1 {
		t.Fatalf("expected wrap to last item, got index %d", m.Index())
	}

	m.CursorDown()
	if m.Index() != 0 {
		t.Fatalf("expected wrap to first item, got index %d", m.Index())
	}

	m.NextPage()
	m.PrevPage()
}

func TestFilterStatesAndHideFlow(t *testing.T) {
	m := newModel([]itemlist.Item{
		extItem{title: "alpha"},
		extItem{title: "beta", hidden: true},
	})

	m.SetFilterText("alpha")
	if m.FilterState() != itemlist.FilterApplied {
		t.Fatalf("expected filter applied, got %v", m.FilterState())
	}
	if !m.IsFiltered() {
		t.Fatalf("expected filtered state")
	}
	if m.FilterValue() != "alpha" {
		t.Fatalf("unexpected filter value %q", m.FilterValue())
	}

	m.ToggleHide()
	visible := m.VisibleItems()
	if len(visible) != 1 {
		t.Fatalf("expected 1 visible item after hide toggle, got %d", len(visible))
	}

	m.ResetFilter()
	if m.FilterState() != itemlist.Unfiltered {
		t.Fatalf("expected unfiltered state, got %v", m.FilterState())
	}
}

func TestSpinnerBindingsHelpAndView(t *testing.T) {
	m := newModel([]itemlist.Item{extItem{title: "alpha"}})

	if !m.KeyMap.Quit.Enabled() {
		t.Fatalf("expected quit binding enabled by default")
	}
	m.DisableQuitKeybindings()
	if m.KeyMap.Quit.Enabled() {
		t.Fatalf("expected quit binding disabled")
	}

	if cmd := m.StartSpinner(); cmd == nil {
		t.Fatalf("expected start spinner command")
	}
	if cmd := m.ToggleSpinner(); cmd != nil {
		t.Fatalf("expected nil toggle command when stopping spinner")
	}
	m.StopSpinner()

	short := m.ShortHelp()
	full := m.FullHelp()
	if len(short) == 0 || len(full) == 0 {
		t.Fatalf("expected help bindings to be available")
	}

	view := m.View().Content
	if view == "" {
		t.Fatalf("expected non-empty view content")
	}
	if !strings.Contains(view, "alpha") {
		t.Fatalf("expected item title in view, got %q", view)
	}
}

func TestFiltersAndStateString(t *testing.T) {
	ranks := itemlist.UnsortedFilter("a", []string{"ba", "ab"})
	if len(ranks) == 0 {
		t.Fatalf("expected at least one rank")
	}

	if itemlist.Unfiltered.String() == "" || itemlist.Filtering.String() == "" ||
		itemlist.FilterApplied.String() == "" {
		t.Fatalf("expected non-empty filter state strings")
	}
}
