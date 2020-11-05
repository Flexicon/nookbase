package categories

import (
	"fmt"
	"strconv"
	"strings"
)

func fishFromRow(row []interface{}) CategoryItem {
	number, _ := strconv.ParseInt(fmt.Sprintf("%v", row[0]), 10, 64)
	sell, _ := strconv.ParseInt(fmt.Sprintf("%v", row[5]), 10, 64)
	totalCatches, _ := strconv.ParseInt(fmt.Sprintf("%v", row[10]), 10, 64)
	rain := strings.Contains(row[6].(string), "rain")

	return CategoryItem{
		Name:                 row[1].(string),
		Number:               number,
		IconImage:            imageFromCell(row[2].(string)),
		CritterpediaImage:    imageFromCell(row[3].(string)),
		FurnitureImage:       imageFromCell(row[4].(string)),
		Sell:                 sell,
		WhereHow:             row[6].(string),
		RainyDays:            &rain,
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
