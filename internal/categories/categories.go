package categories

import (
	"log"
	"strings"

	"github.com/flexicon/nookbase/internal/categories/names"
)

// BuildCategoryFromRow and return it
func BuildCategoryFromRow(category string, row []interface{}) CategoryItem {
	// TODO: map remaining categories
	switch category {
	case names.Insects:
		return insectFromRow(row)
	case names.Fish:
		return fishFromRow(row)
	case names.SeaCreatures:
		return seaCreatureFromRow(row)
	default:
		return defaultFromRow(row)
	}
}

// CategoryItem - generic item with all potential fields
type CategoryItem struct {
	Name                 string              `json:"name"`
	Image                string              `json:"image,omitempty"`
	Number               int64               `json:"number,omitempty"`
	IconImage            string              `json:"icon_image,omitempty"`
	CritterpediaImage    string              `json:"critterpedia_image,omitempty"`
	FurnitureImage       string              `json:"furniture_image,omitempty"`
	Sell                 int64               `json:"sell,omitempty"`
	WhereHow             string              `json:"where_how,omitempty"`
	Weather              string              `json:"weather,omitempty"`
	MovementSpeed        string              `json:"movement_speed,omitempty"`
	Shadow               string              `json:"shadow,omitempty"`
	CatchDifficulty      string              `json:"catch_difficulty,omitempty"`
	Vision               string              `json:"vision,omitempty"`
	TotalCatchesToUnlock int64               `json:"total_catches_to_unlock,omitempty"`
	SpawnRates           string              `json:"spawn_rates,omitempty"`
	NHAvailability       *YearlyAvailability `json:"nh_availability,omitempty"`
	SHAvailability       *YearlyAvailability `json:"sh_availability,omitempty"`
	Size                 string              `json:"size,omitempty"`
	Surface              string              `json:"surface,omitempty"`
	Description          string              `json:"description,omitempty"`
	Catchphrase          string              `json:"catchphrase,omitempty"`
	UniqueID             string              `json:"unique_id,omitempty"`
}

func defaultFromRow(row []interface{}) CategoryItem {
	return CategoryItem{
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

func yearlyAvailabilityFromMonthCols(months []interface{}) *YearlyAvailability {
	avail := make(YearlyAvailability)

	for i, m := range months {
		avail[i+1] = m.(string)
	}

	return &avail
}
