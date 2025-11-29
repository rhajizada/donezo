package navigation

import "testing"

func TestViewOrdering(t *testing.T) {
	expected := []View{
		ViewBoards,
		ViewTags,
		ViewItemsByBoard,
		ViewItemsByTag,
	}

	for i, v := range expected {
		if int(v) != i {
			t.Fatalf("expected view %v to equal %d", v, i)
		}
	}
}
