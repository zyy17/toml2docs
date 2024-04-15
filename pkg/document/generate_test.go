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
		{
			input:    "./testdata/array-table/input.toml",
			expected: "./testdata/array-table/expected.md",
		},
		{
			input:    "./testdata/comment-metadata/input.toml",
			expected: "./testdata/comment-metadata/expected.md",
		},
	}

	for _, tt := range tests {
		got, err := GenerateMarkdownFromFile(tt.input, nil)
		if err != nil {
			t.Fatalf("failed to generate markdown from file '%s': %v", tt.input, err)
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

func TestGenerateMarkdownFromTemplate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "./testdata/templates/input.template",
			expected: "./testdata/templates/expected.md",
		},
	}

	for _, tt := range tests {
		got, err := GenerateMarkdownFromTemplate(tt.input, nil)
		if err != nil {
			t.Fatalf("failed to generate markdown from file '%s': %v", tt.input, err)
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
