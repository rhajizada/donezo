package itemlist

import (
	"io"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
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
	tests := []struct {
		name       string
		filter     string
		wantState  FilterState
		wantTitles []string
	}{
		{name: "filter narrows visible items", filter: "beta", wantState: FilterApplied, wantTitles: []string{"beta"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newModel([]Item{stubItem{title: "alpha"}, stubItem{title: "beta"}})
			m.SetFilterText(tt.filter)
			assert.Equal(t, tt.wantState, m.FilterState())
			assert.Equal(t, tt.filter, m.FilterValue())
			visible := m.VisibleItems()
			assert.Len(t, visible, len(tt.wantTitles))
			for i, want := range tt.wantTitles {
				assert.Equal(t, want, visible[i].(stubItem).title)
			}
		})
	}
}

func TestToggleHideRespectsHideValue(t *testing.T) {
	tests := []struct {
		name      string
		wantCount []int
	}{
		{name: "toggle hide filters hidden items then restores them", wantCount: []int{2, 1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newModel([]Item{stubItem{title: "keep"}, stubItem{title: "hide", hidden: true}})
			assert.Len(t, m.VisibleItems(), tt.wantCount[0])
			m.ToggleHide()
			visible := m.VisibleItems()
			assert.Len(t, visible, tt.wantCount[1])
			assert.Equal(t, "keep", visible[0].(stubItem).title)
			m.ToggleHide()
			assert.Len(t, m.VisibleItems(), tt.wantCount[2])
		})
	}
}

func TestFilteringFlowEnablesBindings(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "filtering flow enables accept and applies filter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newModel([]Item{stubItem{title: "alpha"}, stubItem{title: "beta"}})
			model, _ := m.Update(tea.KeyPressMsg{Code: '/', Text: "/"})
			m = model
			assert.Equal(t, Filtering, m.FilterState())
			assert.False(t, m.KeyMap.AcceptWhileFiltering.Enabled())

			model, cmd := m.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
			m = model
			if cmd != nil {
				if msg := cmd(); msg != nil {
					model, _ = m.Update(msg)
					m = model
				}
			}
			assert.True(t, m.KeyMap.AcceptWhileFiltering.Enabled())

			model, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			m = model
			assert.Equal(t, FilterApplied, m.FilterState())
		})
	}
}
