package main

import (
	"fmt"
	"sort"

	"github.com/flexicon/nookbase/internal/categories"
	"github.com/flexicon/nookbase/internal/categories/names"
	"google.golang.org/api/sheets/v4"
)

// ErrorResponse JSON format
type ErrorResponse struct {
	Error string                 `json:"error"`
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// InvalidCategoryError for the given CategoriesMap
func InvalidCategoryError(c CategoriesMap) ErrorResponse {
	return ErrorResponse{Error: "invalid category", Extra: map[string]interface{}{"categories": c.List()}}
}

// CategoryRow represents a placeholder for the category row response
type CategoryRow interface{}

// Category in dataset
type Category interface {
	// Name returns the category name
	Name() string
	// NameColumn returns a string, referring to the column range in which an category's name can be found
	NameColumn() string
	// MapValueRanges to response models
	MapValueRanges(ranges []*sheets.ValueRange) []CategoryRow
}

// CategoriesMap for a given resource
type CategoriesMap map[string]Category

// List out all available categories in map
func (m CategoriesMap) List() []string {
	var cats []string
	for c := range m {
		cats = append(cats, c)
	}
	sort.Strings(cats)

	return cats
}

// HasKey checks if a key exists in CategoriesMap
func (m CategoriesMap) HasKey(key string) bool {
	_, hasKey := m[key]
	return hasKey
}

// Get a category for the given name
func (m CategoriesMap) Get(categoryName string) Category {
	if !m.HasKey(categoryName) {
		return nil
	}

	category := m[categoryName]
	if category == nil {
		category = newDefaultCategory(categoryName)
	}

	return category
}

// DefaultCategory represents a default category where no special cases are needed to map or query sheet data
type defaultCategory struct {
	name string
}

func newDefaultCategory(name string) defaultCategory {
	return defaultCategory{name: name}
}

func (c defaultCategory) Name() string {
	return c.name
}

func (c defaultCategory) NameColumn() string {
	return fmt.Sprintf("%s!A2:A", c.name)
}

func (c defaultCategory) MapValue(vRange *sheets.ValueRange) CategoryRow {
	return categories.BuildCategoryFromRow(c.Name(), vRange.Values[0])
}

func (c defaultCategory) MapValueRanges(ranges []*sheets.ValueRange) []CategoryRow {
	rows := make([]CategoryRow, 0)

	for _, value := range ranges {
		rows = append(rows, c.MapValue(value))
	}

	return rows
}

// insectsCategory - https://docs.google.com/spreadsheets/d/13d_LAJPlxMa_DubPTuirkIV4DERBMXbrWQsmSh8ReK4/edit#gid=1638053417
type insectsCategory struct {
	defaultCategory
}

func newInsectsCategory() insectsCategory {
	return insectsCategory{defaultCategory: newDefaultCategory(names.Insects)}
}

func (c insectsCategory) NameColumn() string {
	return fmt.Sprintf("%s!B2:B", c.name)
}

// fishCategory - https://docs.google.com/spreadsheets/d/13d_LAJPlxMa_DubPTuirkIV4DERBMXbrWQsmSh8ReK4/edit#gid=1111506211
type fishCategory struct {
	defaultCategory
}

func newFishCategory() fishCategory {
	return fishCategory{defaultCategory: newDefaultCategory(names.Fish)}
}

func (c fishCategory) NameColumn() string {
	return fmt.Sprintf("%s!B2:B", c.name)
}

// seaCreaturesCategory - https://docs.google.com/spreadsheets/d/13d_LAJPlxMa_DubPTuirkIV4DERBMXbrWQsmSh8ReK4/edit#gid=60735325
type seaCreaturesCategory struct {
	defaultCategory
}

func newSeaCreaturesCategory() seaCreaturesCategory {
	return seaCreaturesCategory{defaultCategory: newDefaultCategory(names.SeaCreatures)}
}

func (c seaCreaturesCategory) NameColumn() string {
	return fmt.Sprintf("%s!B2:B", c.name)
}
