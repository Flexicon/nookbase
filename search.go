package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/flexicon/nookbase/internal/categories/names"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"google.golang.org/api/sheets/v4"
)

const (
	// SearchBatchLimit - the maximum amount of ranges that we want to query in a single request to Google Sheets API
	SearchBatchLimit = 400
)

var (
	// SearchCategories that are available - nil values use the defaultCategory
	SearchCategories = CategoriesMap{
		names.Housewares:    nil,
		names.Miscellaneous: nil,
		names.WallMounted:   nil,
		names.Wallpaper:     nil,
		names.Floors:        nil,
		names.Rugs:          nil,
		names.Photos:        nil,
		names.Posters:       nil,
		names.Tools:         nil,
		names.Fencing:       nil,
		names.Tops:          nil,
		names.Bottoms:       nil,
		names.DressUp:       nil,
		names.Headwear:      nil,
		names.Accessories:   nil,
		names.Socks:         nil,
		names.Shoes:         nil,
		names.Bags:          nil,
		names.Umbrellas:     nil,
		names.ClothingOther: nil,
		names.Music:         nil,
		names.Insects:       newInsectsCategory(),
		names.Fish:          newFishCategory(),
		names.SeaCreatures:  newSeaCreaturesCategory(),
		names.Fossils:       nil,
		names.Art:           nil,
		names.Recipes:       nil,
		names.Other:         nil,
		names.Construction:  nil,
		names.Achievements:  nil,
		names.Villagers:     nil,
		names.SpecialNPCs:   nil,
		names.Reactions:     nil,
		names.MessageCards:  nil,
	}
)

// SearchHandler - performs a lookup for items in the given category
//
// GET /search/:category
func SearchHandler(service *sheets.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		category := SearchCategories.Get(sanitizeCategory(c.Param("category")))
		if category == nil {
			return c.JSON(http.StatusBadRequest, InvalidCategoryError(SearchCategories))
		}

		q := strings.ToLower(c.QueryParam("q"))
		if q == "" {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "parameter q is required"})
		}

		nameResults, err := getNamesForCategory(service, category)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}

		matchIndexes := findMatchesInValues(nameResults.Values, q)

		rowResults, err := getRowsByIndexes(service, category.Name(), matchIndexes)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, category.MapValueRanges(rowResults))
	}
}

func findMatchesInValues(values [][]interface{}, query string) []int {
	var matchIndexes []int
	for i, cells := range values {
		if len(cells) != 0 && strings.Contains(strings.ToLower(cells[0].(string)), query) {
			matchIndexes = append(matchIndexes, i+2)
		}
	}

	return matchIndexes
}

func getNamesForCategory(service *sheets.Service, category Category) (*sheets.ValueRange, error) {
	// Lookup all item names in category - first checking cache, then making a request if is a misss
	var results *sheets.ValueRange
	cachedRes, hit := Cache[category.Name()]

	if hit {
		results = cachedRes.(*sheets.ValueRange)
	} else {
		apiRes, err := service.Spreadsheets.Values.Get(SpreadsheetID, category.NameColumn()).Do()
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve names")
		}
		Cache[fmt.Sprintf("%s", category.Name())] = apiRes
		results = apiRes
	}

	return results, nil
}

func getRowsByIndexes(service *sheets.Service, category string, indexes []int) ([]*sheets.ValueRange, error) {
	var rows []*sheets.ValueRange
	var remainingRanges []string

	for _, rangeStr := range indexesToRanges(category, indexes) {
		if cachedRow, hit := Cache[rangeStr]; hit {
			rows = append(rows, cachedRow.(*sheets.ValueRange))
		} else {
			remainingRanges = append(remainingRanges, rangeStr)
		}
	}

	if len(remainingRanges) != 0 {
		for i := 0; len(remainingRanges)-i > 1; i += SearchBatchLimit {
			upperLimit := i + SearchBatchLimit
			if upperLimit > len(remainingRanges) {
				upperLimit = i + (len(remainingRanges) % SearchBatchLimit)
			}

			call := service.Spreadsheets.Values.BatchGet(SpreadsheetID)
			results, err := call.Ranges(remainingRanges[i:upperLimit]...).ValueRenderOption("FORMULA").Do()
			if err != nil {
				return nil, errors.Wrap(err, "failed to retrieve matched rows")
			}

			for j, row := range results.ValueRanges {
				Cache[remainingRanges[j+i]] = row
				rows = append(rows, row)
			}
		}
	}

	return rows, nil
}

func indexesToRanges(sheet string, indexes []int) []string {
	var ranges []string
	for _, i := range indexes {
		ranges = append(ranges, fmt.Sprintf("%s!%d:%d", sheet, i, i))
	}

	return ranges
}

func sanitizeCategory(category string) string {
	return strings.ToLower(strings.ReplaceAll(category, "_", " "))
}
