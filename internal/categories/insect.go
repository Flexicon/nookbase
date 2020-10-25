package categories

import (
	"fmt"
	"strconv"
)

// Insect represents a listing for the category of the same name
type Insect struct {
	Number               int64              `json:"number"`
	Name                 string             `json:"name"`
	IconImage            string             `json:"icon_image"`
	CritterpediaImage    string             `json:"critterpedia_image"`
	FurnitureImage       string             `json:"furniture_image"`
	Sell                 int64              `json:"sell"`
	WhereHow             string             `json:"where_how"`
	Weather              string             `json:"weather"`
	TotalCatchesToUnlock int64              `json:"total_catches_to_unlock"`
	SpawnRates           string             `json:"spawn_rates"`
	NHAvailability       YearlyAvailability `json:"nh_availability"`
	SHAvailability       YearlyAvailability `json:"sh_availability"`
	Size                 string             `json:"size"`
	Surface              string             `json:"surface"`
	Description          string             `json:"description"`
	Catchphrase          string             `json:"catchphrase"`
	UniqueID             string             `json:"unique_id"`
}

func insectFromRow(row []interface{}) Insect {
	number, _ := strconv.ParseInt(fmt.Sprintf("%v", row[0]), 10, 64)
	sell, _ := strconv.ParseInt(fmt.Sprintf("%v", row[5]), 10, 64)
	totalCatches, _ := strconv.ParseInt(fmt.Sprintf("%v", row[8]), 10, 64)

	return Insect{
		Name:                 row[1].(string),
		Number:               number,
		IconImage:            imageFromCell(row[2].(string)),
		CritterpediaImage:    imageFromCell(row[3].(string)),
		FurnitureImage:       imageFromCell(row[4].(string)),
		Sell:                 sell,
		WhereHow:             row[6].(string),
		Weather:              row[7].(string),
		TotalCatchesToUnlock: totalCatches,
		SpawnRates:           row[9].(string),
		NHAvailability:       yearlyAvailabilityFromMonthCols(row[10:22]),
		SHAvailability:       yearlyAvailabilityFromMonthCols(row[22:34]),
		Size:                 row[34].(string),
		Surface:              row[35].(string),
		Description:          row[36].(string),
		Catchphrase:          row[37].(string),
		UniqueID:             row[46].(string),
	}
}
