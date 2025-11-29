package itemlist

import (
	"io"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type stubItem struct {
	title  string
	hidden bool
}

func (s stubItem) FilterValue() string { return s.title }
func (s stubItem) HideValue() bool     { return s.hidden }

type stubDelegate struct {
	height  int
	spacing int
}

func (d stubDelegate) Render(io.Writer, Model, int, Item) {}
func (d stubDelegate) Height() int {
	if d.height == 0 {
		return 1
	}
	return d.height
}
func (d stubDelegate) Spacing() int { return d.spacing }
func (d stubDelegate) Update(tea.Msg, *Model) tea.Cmd {
	return nil
}

func newModel(items []Item) Model {
	return New(items, stubDelegate{height: 1, spacing: 0}, 80, 20)
}

func TestSetFilterTextAppliesMatches(t *testing.T) {
	items := []Item{
		stubItem{title: "alpha"},
		stubItem{title: "beta"},
	}
	m := newModel(items)

	m.SetFilterText("beta")

	if m.FilterState() != FilterApplied {
		t.Fatalf("expected filter applied state, got %v", m.FilterState())
	}
	if m.FilterValue() != "beta" {
		t.Fatalf("unexpected filter value %q", m.FilterValue())
	}
	visible := m.VisibleItems()
	if len(visible) != 1 {
		t.Fatalf("expected 1 visible item after filter, got %d", len(visible))
	}
	if visible[0].(stubItem).title != "beta" {
		t.Fatalf("unexpected filtered item %v", visible[0])
	}
}

func TestToggleHideRespectsHideValue(t *testing.T) {
	items := []Item{
		stubItem{title: "keep"},
		stubItem{title: "hide", hidden: true},
	}
	m := newModel(items)

	if len(m.VisibleItems()) != 2 {
		t.Fatalf("expected 2 items initially, got %d", len(m.VisibleItems()))
	}

	m.ToggleHide()
	visible := m.VisibleItems()
	if len(visible) != 1 {
		t.Fatalf("expected 1 visible item after hiding, got %d", len(visible))
	}
	if visible[0].(stubItem).title != "keep" {
		t.Fatalf("unexpected item kept visible: %v", visible[0])
	}

	m.ToggleHide()
	if len(m.VisibleItems()) != 2 {
		t.Fatalf("expected hidden items restored, got %d", len(m.VisibleItems()))
	}
}

func TestFilteringFlowEnablesBindings(t *testing.T) {
	items := []Item{
		stubItem{title: "alpha"},
		stubItem{title: "beta"},
	}
	m := newModel(items)

	// Enter filtering mode.
	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = model
	if m.FilterState() != Filtering {
		t.Fatalf("expected filtering state, got %v", m.FilterState())
	}
	if m.KeyMap.AcceptWhileFiltering.Enabled() {
		t.Fatalf("accept should be disabled with empty filter")
	}

	// Type a filter value.
	model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = model
	if cmd != nil {
		if msg := cmd(); msg != nil {
			model, _ = m.Update(msg)
			m = model
		}
	}
	if !m.KeyMap.AcceptWhileFiltering.Enabled() {
		t.Fatalf("accept should enable after typing filter text")
	}

	// Accept the filter.
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = model
	if m.FilterState() != FilterApplied {
		t.Fatalf("expected filter applied after accept, got %v", m.FilterState())
	}
}
