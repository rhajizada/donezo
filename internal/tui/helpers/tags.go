package helpers

import (
	"errors"
	"strings"
)

const TagsSeparator = ","

func ExtractTags(input string) ([]string, error) {
	tags := strings.Split(input, TagsSeparator)

	if len(tags) == 1 && tags[0] == "" {
		tags = make([]string, 0)
	}

	for idx, tag := range tags {
		sanitized := strings.TrimSpace(tag)
		if len(sanitized) == 0 {
			return nil, errors.New("tag must not be empty string")
		}

		tags[idx] = sanitized
	}

	return tags, nil
}
