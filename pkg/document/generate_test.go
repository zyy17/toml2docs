package document

import (
	"os"
	"testing"
)

func TestGenerateMarkdown(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "./testdata/basic/input.toml",
			expected: "./testdata/basic/expected.md",
		},
		{
			input:    "./testdata/nested/input.toml",
			expected: "./testdata/nested/expected.md",
		},
	}

	for _, tt := range tests {
		data, err := os.ReadFile(tt.input)
		if err != nil {
			t.Fatalf("failed to read the input file: %v", err)
		}

		got, err := GenerateMarkdown(data, nil)
		if err != nil {
			t.Fatalf("failed to generate markdown: %v", err)
		}

		expected, err := os.ReadFile(tt.expected)
		if err != nil {
			t.Fatalf("failed to read the expected file: %v", err)
		}

		if got != string(expected) {
			t.Errorf("expected: %s, got: %s", string(expected), got)
		}
	}
}
