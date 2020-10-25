package categories

import (
	"log"
	"strings"

	"github.com/flexicon/nookbase/internal/categories/names"
)

// BuildCategoryFromRow and return it
func BuildCategoryFromRow(category string, row []interface{}) interface{} {
	// TODO: map remaining categories
	switch category {
	case names.Fish:
		return fishFromRow(row)
	case names.SeaCreatures:
		return seaCreatureFromRow(row)
	default:
		return defaultFromRow(row)
	}
}

// Default represents a default listing with common fields
type Default struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func defaultFromRow(row []interface{}) Default {
	return Default{
		Name:  row[0].(string),
		Image: imageFromCell(row[1].(string)),
	}
}

func imageFromCell(cell string) string {
	parts := strings.Split(cell, `"`)

	if len(parts) == 1 {
		log.Printf("image cell in unexpected format: %s", cell)
		return ""
	}
	return parts[1]
}

// YearlyAvailability map of items per month
type YearlyAvailability map[int]string

func yearlyAvailabilityFromMonthCols(months []interface{}) YearlyAvailability {
	avail := make(YearlyAvailability)

	for i, m := range months {
		avail[i+1] = m.(string)
	}

	return avail
}
