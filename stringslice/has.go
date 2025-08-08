package stringslice

import "slices"

// Has returns true if the needle is in the haystack (case-sensitive)
// Deprecated: use slices.Contains instead
func Has(haystack []string, needle string) bool {
	return slices.Contains(haystack, needle)
}
