package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitTranslations(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected []string
	}{
		{
			name:     "Empty",
			in:       "",
			expected: []string{},
		},
		{
			name:     "One entry",
			in:       "word",
			expected: []string{"word"},
		},
		{
			name:     "Several entries",
			in:       "word,another,something",
			expected: []string{"word", "another", "something"},
		},
		{
			name:     "Several entries with spaces",
			in:       "word,   another, something",
			expected: []string{"word", "another", "something"},
		},
		{
			name:     "Several entries with empty",
			in:       "word,,  ,  another,     , something",
			expected: []string{"word", "another", "something"},
		},
		{
			name:     "Several entries, all empty",
			in:       ",,  ,  ,     ,",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitTranslations(tt.in)

			assert.Equal(t, tt.expected, result)
		})
	}

}
