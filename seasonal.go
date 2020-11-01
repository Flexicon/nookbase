package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/flexicon/nookbase/internal/categories/names"
	"github.com/flexicon/nookbase/internal/mapping"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"google.golang.org/api/sheets/v4"
)

var (
	// SeasonalCategories that are available
	SeasonalCategories = CategoriesMap{
		names.Insects:      newInsectsCategory(),
		names.Fish:         newFishCategory(),
		names.SeaCreatures: newSeaCreaturesCategory(),
	}
)

// SeasonalHandler - performs a lookup for items that are currently in season for the given hemisphere and category
//
// GET /seasonal/:hemisphere/:category
func SeasonalHandler(service *sheets.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Validate inputs
		hemisphere := Hemisphere(c.Param("hemisphere"))
		if !hemisphere.IsValid() {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid hemisphere - should be either northern or southern"})
		}

		category := SeasonalCategories.Get(mapping.NormalizeCategory(c.Param("category")))
		if category == nil {
			return c.JSON(http.StatusBadRequest, InvalidCategoryError(SeasonalCategories))
		}

		// Get all availabilities for all items in category
		availability, err := getAvailabilityForCategory(service, category, hemisphere)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}
		// Narrow down to only ones available in the current season
		availableIndexes := findAvailableForMonthInValues(availability.Values, int(time.Now().Month()))

		// Retrieve full rows for matched indexes
		rowResults, err := getRowsByIndexes(service, category.Name(), availableIndexes)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, category.MapValueRanges(rowResults))
	}
}

func getAvailabilityForCategory(service *sheets.Service, category Category, hemisphere Hemisphere) (*sheets.ValueRange, error) {
	var results *sheets.ValueRange
	cachedRes, hit := Cache[fmt.Sprintf("%s-%s-Availability", category.Name(), hemisphere)]

	if hit {
		results = cachedRes.(*sheets.ValueRange)
	} else {
		apiRes, err := service.Spreadsheets.Values.Get(SpreadsheetID, category.AvailabilityRange(hemisphere)).Do()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to retrieve %s availability", hemisphere))
		}
		Cache[fmt.Sprintf("%s-%s-Availability", category.Name(), hemisphere)] = apiRes
		results = apiRes
	}

	return results, nil
}

func findAvailableForMonthInValues(values [][]interface{}, month int) []int {
	var matchIndexes []int
	for i, cells := range values {
		if len(cells) >= month && strings.ToUpper(cells[month-1].(string)) != "NA" {
			matchIndexes = append(matchIndexes, i+2)
		}
	}

	return matchIndexes
}
