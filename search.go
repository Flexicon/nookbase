package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/flexicon/nookbase/internal/categories/names"
	"github.com/flexicon/nookbase/internal/mapping"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"google.golang.org/api/sheets/v4"
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
// @Summary Performs a search for items in the given category
// @Tags search
//
// @Accept  json
// @Produce  json
//
// @Param q query string true "The search query"
// @Param category path string true "Item category to search through (eg: fish, insects, etc.)"
//
// @Success 200 {object} []categories.CategoryItem
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
//
// @Router /search/{category} [get]
func SearchHandler(service *sheets.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		category := SearchCategories.Get(mapping.NormalizeCategory(c.Param("category")))
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

		matchIndexes := findQueryMatchesInValues(nameResults.Values, q)

		rowResults, err := getRowsByIndexes(service, category.Name(), matchIndexes)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, category.MapValueRanges(rowResults))
	}
}

func findQueryMatchesInValues(values [][]interface{}, query string) []int {
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
