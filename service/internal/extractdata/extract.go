// Package extractdata provides the logic for extracting data from input.
package extractdata

import (
	"context"
	"fmt"

	"github.com/mnabbasabadi/relex_convertor/service/shared/domain"
)

var _ Extractor = new(controller)

type (
	// Extractor is an interface for extracting data from a source
	Extractor interface {
		Extract(context.Context, [][]string) ([]domain.Data, error)
	}
	controller struct {
	}
)

// New ...
func New() Extractor {
	return &controller{}
}

// Extract ...
func (c controller) Extract(_ context.Context, records [][]string) ([]domain.Data, error) {
	ln := len(records)
	if ln == 0 {
		return nil, fmt.Errorf("no records found")
	}
	colNum := len(records[0])
	if colNum < 2 || colNum > 4 {
		return nil, fmt.Errorf("invalid number of records")
	}

	// iterate through the records
	level1Index, level2Index, level3Index, itemIDIndex := findIndexes(records[0])

	dataSlice := make([]domain.Data, 0, ln-1)

	for i, record := range records {
		if i == 0 {
			// skip the header
			continue
		}
		data := domain.Data{
			Level1: record[*level1Index],
			ItemID: record[*itemIDIndex],
		}
		if level2Index != nil {
			data.Level2 = record[*level2Index]
		}
		if level3Index != nil {
			data.Level3 = record[*level3Index]
		}

		if err := validate(data.Level1, data.ItemID, data.Level2, data.Level3); err != nil {
			return nil, err
		}

		dataSlice = append(dataSlice, data)

	}

	return dataSlice, nil
}

func findIndexes(fields []string) (level1Index, level2Index, level3Index, itemIDIndex *int) {
	for i, field := range fields {
		index := i
		switch field {
		case "level_1":
			level1Index = &index
		case "level_2":
			level2Index = &index
		case "level_3":
			level3Index = &index
		case "item_id":
			itemIDIndex = &index
		}
	}
	return
}

func validate(level1 string, itemID string, level2 string, level3 string) error {
	if level1 == "" {
		return fmt.Errorf("missing Level1 value")
	}
	if itemID == "" {
		return fmt.Errorf("missing ItemID value")
	}

	if level2 == "" && level3 != "" {
		return fmt.Errorf("invalid structure: Missing Level3 value for Level2 %s", level2)
	}
	return nil
}
