package utils

import "testing"

func TestNormalizeSourceURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic URL without trailing slash",
			input:    "https://api.scanoss.com/file_contents",
			expected: "https://api.scanoss.com/file_contents/", // Current function adds a slash
		},
		{
			name:     "URL with trailing slash",
			input:    "https://api.scanoss.com/file_contents/",
			expected: "https://api.scanoss.com/file_contents/", // Current function leaves this unchanged
		},
		{
			name:     "URL with leading spaces",
			input:    "  https://api.scanoss.com/file_contents",
			expected: "https://api.scanoss.com/file_contents/", // Current function trims spaces and adds slash
		},
		{
			name:     "URL with trailing spaces",
			input:    "https://api.scanoss.com/file_contents  ",
			expected: "https://api.scanoss.com/file_contents/", // Current function trims spaces and adds slash
		},
		{
			name:     "URL with leading and trailing spaces",
			input:    "  https://api.scanoss.com/file_contents  ",
			expected: "https://api.scanoss.com/file_contents/", // Current function trims spaces and adds slash
		},
		{
			name:     "URL with trailing slash and spaces",
			input:    "https://api.scanoss.com/file_contents/  ",
			expected: "https://api.scanoss.com/file_contents/", // Current function trims spaces and keeps slash
		},
		{
			name:     "Multiple trailing slashes",
			input:    "https://api.scanoss.com/file_contents///",
			expected: "https://api.scanoss.com/file_contents/", // Current function normalize the trailing slash
		},
		{
			name:     "Multiple trailing slashes with spaces",
			input:    "    https://api.scanoss.com/file_contents///     ",
			expected: "https://api.scanoss.com/file_contents/", // Current function normalize the trailing slash
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NormalizeSourceURL(tc.input)
			if result != tc.expected {
				t.Errorf("normalizeSourceURL(%q) = %q, expected %q", tc.input, result, tc.expected)
			}

		})
	}
}
