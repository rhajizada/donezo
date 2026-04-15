package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr string
	}{
		{name: "trims spaces", input: "work, personal , errands", want: []string{"work", "personal", "errands"}},
		{name: "empty string yields none", input: "", want: []string{}},
		{name: "single tag", input: "one", want: []string{"one"}},
		{name: "error on empty tag", input: "a,", wantErr: "tag must not be empty string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tags, err := ExtractTags(tt.input)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, tags)
		})
	}
}
