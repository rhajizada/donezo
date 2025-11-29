package helpers

import (
	"reflect"
	"testing"
)

func TestExtractTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		want      []string
		expectErr bool
	}{
		{name: "trims spaces", input: "work, personal , errands", want: []string{"work", "personal", "errands"}},
		{name: "empty string yields none", input: "", want: []string{}},
		{name: "single tag", input: "one", want: []string{"one"}},
		{name: "error on empty tag", input: "a,", expectErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tags, err := ExtractTags(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(tags, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, tags)
			}
		})
	}
}
