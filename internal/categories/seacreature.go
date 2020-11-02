package categories

import (
	"fmt"
	"strconv"
)

func seaCreatureFromRow(row []interface{}) CategoryItem {
	number, _ := strconv.ParseInt(fmt.Sprintf("%v", row[0]), 10, 64)
	sell, _ := strconv.ParseInt(fmt.Sprintf("%v", row[5]), 10, 64)
	totalCatches, _ := strconv.ParseInt(fmt.Sprintf("%v", row[8]), 10, 64)

	return CategoryItem{
		Name:                 row[1].(string),
		Number:               number,
		IconImage:            imageFromCell(row[2].(string)),
		CritterpediaImage:    imageFromCell(row[3].(string)),
		FurnitureImage:       imageFromCell(row[4].(string)),
		Sell:                 sell,
		Shadow:               row[6].(string),
		MovementSpeed:        row[7].(string),
		TotalCatchesToUnlock: totalCatches,
		SpawnRates:           row[9].(string),
		NHAvailability:       yearlyAvailabilityFromMonthCols(row[10:22]),
		SHAvailability:       yearlyAvailabilityFromMonthCols(row[22:34]),
		Size:                 row[34].(string),
		Surface:              row[35].(string),
		Description:          row[36].(string),
		Catchphrase:          row[37].(string),
		UniqueID:             row[49].(string),
	}
}
