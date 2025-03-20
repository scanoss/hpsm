package utils

import "strings"

// NormalizeSourceURL processes a URL to ensure consistent formatting
// Current implementation:
// 1. Trims leading spaces with TrimLeft
// 2. Trims trailing spaces AND slashes with TrimRight
// 3. Adds a single trailing slash to ensure consistent URL format
//
// This approach standardizes URLs by ensuring exactly one trailing slash,
// preventing double-slash issues when concatenating paths.
func NormalizeSourceURL(baseurl string) string {
	leftTrimmedUrl := strings.TrimLeft(baseurl, " ")
	trimmedUrl := strings.TrimRight(leftTrimmedUrl, "/ ")
	return trimmedUrl + "/"
}
