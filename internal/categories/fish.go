package categories

import (
	"fmt"
	"strconv"
)

// Fish represents a listing for the category of the same name
type Fish struct {
	Number               int64              `json:"number"`
	Name                 string             `json:"name"`
	IconImage            string             `json:"icon_image"`
	CritterpediaImage    string             `json:"critterpedia_image"`
	FurnitureImage       string             `json:"furniture_image"`
	Sell                 int64              `json:"sell"`
	WhereHow             string             `json:"where_how"`
	Shadow               string             `json:"shadow"`
	CatchDifficulty      string             `json:"catch_difficulty"`
	Vision               string             `json:"vision"`
	TotalCatchesToUnlock int64              `json:"total_catches_to_unlock"`
	SpawnRates           string             `json:"spawn_rates"`
	NHAvailability       YearlyAvailability `json:"nh_availability"`
	SHAvailability       YearlyAvailability `json:"sh_availability"`

	Size        string `json:"size"`
	Description string `json:"description"`
	Catchphrase string `json:"catchphrase"`
	UniqueID    string `json:"unique_id"`
}

func fishFromRow(row []interface{}) Fish {
	number, _ := strconv.ParseInt(fmt.Sprintf("%v", row[0]), 10, 64)
	sell, _ := strconv.ParseInt(fmt.Sprintf("%v", row[5]), 10, 64)
	totalCatches, _ := strconv.ParseInt(fmt.Sprintf("%v", row[10]), 10, 64)

	return Fish{
		Name:                 row[1].(string),
		Number:               number,
		IconImage:            imageFromCell(row[2].(string)),
		CritterpediaImage:    imageFromCell(row[3].(string)),
		FurnitureImage:       imageFromCell(row[4].(string)),
		Sell:                 sell,
		WhereHow:             row[6].(string),
		Shadow:               row[7].(string),
		CatchDifficulty:      row[8].(string),
		Vision:               row[9].(string),
		TotalCatchesToUnlock: totalCatches,
		SpawnRates:           row[11].(string),
		NHAvailability:       yearlyAvailabilityFromMonthCols(row[12:24]),
		SHAvailability:       yearlyAvailabilityFromMonthCols(row[24:36]),
		Size:                 row[36].(string),
		Description:          row[38].(string),
		Catchphrase:          row[39].(string),
		UniqueID:             row[49].(string),
	}
}
