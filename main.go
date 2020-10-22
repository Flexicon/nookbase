package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	// SpreadsheetID for the Animal Crossing data spreadsheet - https://docs.google.com/spreadsheets/d/13d_LAJPlxMa_DubPTuirkIV4DERBMXbrWQsmSh8ReK4
	SpreadsheetID = "13d_LAJPlxMa_DubPTuirkIV4DERBMXbrWQsmSh8ReK4"
)

var (
	// Cache stores Google Sheets API responses in memory
	Cache = make(map[string]interface{})
)

func main() {
	service, err := sheets.NewService(context.TODO(), option.WithCredentialsFile("./creds.json"))
	if err != nil {
		log.Fatalf("failed to setup sheets service: %v", err)
	}

	e := echo.New()
	e.Debug = true

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{DisableStackAll: true}))
	e.Use(middleware.Secure())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "REQUEST: method=${method}, status=${status}, uri=${uri}, latency=${latency_human}\n",
	}))

	e.GET("/search/:category", searchHandler(service))

	e.Logger.Fatal(e.Start(":3000"))
}

func searchHandler(service *sheets.Service) echo.HandlerFunc {
	type Response struct {
		Msg string `json:"message"`
	}

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	return func(c echo.Context) error {
		category := c.Param("category")
		q := strings.ToLower(c.QueryParam("q"))
		if q == "" {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "parameter q is required"})
		}

		nameResults, err := getNamesForSheet(service, category)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}

		matchIndexes := findMatchesInValues(nameResults.Values, q)

		rowResults, err := getRowsByIndexes(service, category, matchIndexes)
		if err != nil {
			log.Panicln(err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}

		// TODO: properly map results to model structs based on category
		return c.JSON(http.StatusOK, rowResults)
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

func getNamesForSheet(service *sheets.Service, sheet string) (*sheets.ValueRange, error) {
	// Lookup all item names in category - first checking cache, then making a request if is a misss
	var results *sheets.ValueRange
	cachedRes, hit := Cache[sheet]

	if hit {
		results = cachedRes.(*sheets.ValueRange)
	} else {
		apiRes, err := service.Spreadsheets.Values.Get(SpreadsheetID, sheet+"!A2:A").Do()
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve names")
		}
		Cache[fmt.Sprintf("%s", sheet)] = apiRes
		results = apiRes
	}

	return results, nil
}

func getRowsByIndexes(service *sheets.Service, sheet string, indexes []int) ([]*sheets.ValueRange, error) {
	var rows []*sheets.ValueRange
	var remainingRanges []string

	for _, rangeStr := range indexesToRanges(sheet, indexes) {
		if cachedRow, hit := Cache[rangeStr]; hit {
			rows = append(rows, cachedRow.(*sheets.ValueRange))
		} else {
			remainingRanges = append(remainingRanges, rangeStr)
		}
	}

	if len(remainingRanges) != 0 {
		results, err := service.Spreadsheets.Values.BatchGet(SpreadsheetID).Ranges(remainingRanges...).Do()
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve matched rows")
		}

		for i, row := range results.ValueRanges {
			Cache[remainingRanges[i]] = row
			rows = append(rows, row)
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
