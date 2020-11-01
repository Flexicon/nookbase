package main

import (
	"github.com/flexicon/nookbase/internal/mapping"
	"github.com/pkg/errors"
	"google.golang.org/api/sheets/v4"
)

const (
	// RequestBatchLimit - the maximum amount of ranges that we want to query in a single request to Google Sheets API
	RequestBatchLimit = 400
)

func getRowsByIndexes(service *sheets.Service, category string, indexes []int) ([]*sheets.ValueRange, error) {
	var rows []*sheets.ValueRange
	var remainingRanges []string

	for _, rangeStr := range mapping.IndexesToRanges(category, indexes) {
		if cachedRow, hit := Cache[rangeStr]; hit {
			rows = append(rows, cachedRow.(*sheets.ValueRange))
		} else {
			remainingRanges = append(remainingRanges, rangeStr)
		}
	}

	if len(remainingRanges) != 0 {
		for i := 0; len(remainingRanges)-i > 1; i += RequestBatchLimit {
			upperLimit := i + RequestBatchLimit
			if upperLimit > len(remainingRanges) {
				upperLimit = i + (len(remainingRanges) % RequestBatchLimit)
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
