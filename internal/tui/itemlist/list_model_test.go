package itemlist_test

import (
	"fmt"
	"io"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	tests := []struct{ name string }{{name: "public setters mutate model state"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newModel([]itemlist.Item{extItem{title: "a"}})
			m.SetShowTitle(false)
			m.SetShowFilter(false)
			m.SetShowStatusBar(false)
			m.SetShowPagination(false)
			m.SetShowHelp(false)
			m.SetFilteringEnabled(true)
			m.SetStatusBarItemName("task", "tasks")
			m.SetSize(70, 15)

			assert.False(t, m.ShowTitle())
			assert.False(t, m.ShowFilter())
			assert.False(t, m.ShowStatusBar())
			assert.False(t, m.ShowPagination())
			assert.False(t, m.ShowHelp())
			assert.True(t, m.FilteringEnabled())
			singular, plural := m.StatusBarItemName()
			assert.Equal(t, "task", singular)
			assert.Equal(t, "tasks", plural)
			assert.Equal(t, 70, m.Width())
			assert.Equal(t, 15, m.Height())
		})
	}
}

func TestItemMutationAndSelection(t *testing.T) {
	tests := []struct{ name string }{{name: "item mutations preserve selection invariants"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newModel([]itemlist.Item{extItem{title: "a"}, extItem{title: "b"}})
			m.InsertItem(1, extItem{title: "inserted"})
			assert.Len(t, m.Items(), 3)
			m.SetItem(0, extItem{title: "first"})
			m.Select(1)
			assert.Equal(t, 1, m.Index())
			assert.Equal(t, 1, m.GlobalIndex())
			m.RemoveItem(1)
			assert.Len(t, m.Items(), 2)
			m.ResetSelected()
			assert.NotNil(t, m.SelectedItem())
		})
	}
}

func TestCursorPagingAndInfiniteScrolling(t *testing.T) {
	tests := []struct{ name string }{{name: "cursor wraps with infinite scrolling"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := []itemlist.Item{extItem{title: "1"}, extItem{title: "2"}, extItem{title: "3"}}
			m := itemlist.New(items, extDelegate{height: 1, spacing: 0}, 40, 2)
			m.InfiniteScrolling = true
			m.CursorUp()
			assert.Equal(t, len(items)-1, m.Index())
			m.CursorDown()
			assert.Equal(t, 0, m.Index())
			m.NextPage()
			m.PrevPage()
		})
	}
}

func TestFilterStatesAndHideFlow(t *testing.T) {
	tests := []struct{ name string }{{name: "filter state and hide flow update model"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newModel([]itemlist.Item{extItem{title: "alpha"}, extItem{title: "beta", hidden: true}})
			m.SetFilterText("alpha")
			assert.Equal(t, itemlist.FilterApplied, m.FilterState())
			assert.True(t, m.IsFiltered())
			assert.Equal(t, "alpha", m.FilterValue())
			m.ToggleHide()
			assert.Len(t, m.VisibleItems(), 1)
			m.ResetFilter()
			assert.Equal(t, itemlist.Unfiltered, m.FilterState())
		})
	}
}

func TestSpinnerBindingsHelpAndView(t *testing.T) {
	tests := []struct{ name string }{{name: "spinner and help APIs remain usable"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newModel([]itemlist.Item{extItem{title: "alpha"}})
			assert.True(t, m.KeyMap.Quit.Enabled())
			m.DisableQuitKeybindings()
			assert.False(t, m.KeyMap.Quit.Enabled())
			require.NotNil(t, m.StartSpinner())
			assert.Nil(t, m.ToggleSpinner())
			m.StopSpinner()
			assert.NotEmpty(t, m.ShortHelp())
			assert.NotEmpty(t, m.FullHelp())
			view := m.View().Content
			assert.NotEmpty(t, view)
			assert.Contains(t, view, "alpha")
		})
	}
}

func TestFiltersAndStateString(t *testing.T) {
	tests := []struct {
		name  string
		state itemlist.FilterState
	}{
		{name: "unfiltered string", state: itemlist.Unfiltered},
		{name: "filtering string", state: itemlist.Filtering},
		{name: "filter applied string", state: itemlist.FilterApplied},
	}

	ranks := itemlist.UnsortedFilter("a", []string{"ba", "ab"})
	assert.NotEmpty(t, ranks)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.state.String())
		})
	}
}
