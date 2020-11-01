package mapping

import (
	"fmt"
	"strings"
)

// IndexesToRanges maps a given slice of indexes in a sheet to a slice of sheet ranges
func IndexesToRanges(sheet string, indexes []int) []string {
	var ranges []string
	for _, i := range indexes {
		ranges = append(ranges, fmt.Sprintf("%s!%d:%d", sheet, i, i))
	}

	return ranges
}

// NormalizeCategory to match category constants
func NormalizeCategory(category string) string {
	return strings.ToLower(strings.ReplaceAll(category, "_", " "))
}
