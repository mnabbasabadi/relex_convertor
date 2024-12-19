// Package converter provides the logic for converting a CSV file into a JSON tree.
package converter

import (
	"context"

	"github.com/mnabbasabadi/relex_convertor/service/shared/domain"
)

var _ Converter = (*convert)(nil)

type (
	// Converter is an interface for converting data into a tree.
	Converter interface {
		BuildTree(context.Context, []domain.Data) (domain.Node, error)
	}
	convert struct {
	}
)

// New ...
func New() Converter {
	return &convert{}
}

// BuildTree builds a tree from the data records.
func (c convert) BuildTree(_ context.Context, dataRecords []domain.Data) (domain.Node, error) {
	root := domain.Node{
		Item:     false,
		Children: make(map[string]domain.Node),
	}

	for _, record := range dataRecords {
		level1 := record.Level1
		level2 := record.Level2
		level3 := record.Level3
		itemID := record.ItemID

		if _, exists := root.Children[level1]; !exists {
			root.Children[level1] = domain.Node{
				Item:     false,
				Children: make(map[string]domain.Node),
			}
		}

		if level2 != "" {
			if _, exists := root.Children[level1].Children[level2]; !exists {
				root.Children[level1].Children[level2] = domain.Node{
					Item:     false,
					Children: make(map[string]domain.Node),
				}
			}
		} else {
			root.Children[level1].Children[itemID] = domain.Node{
				Item:     true,
				Children: nil,
			}
			continue
		}
		if level3 != "" {
			if _, exists := root.Children[level1].Children[level2].Children[level3]; !exists {
				root.Children[level1].Children[level2].Children[level3] = domain.Node{
					Item:     false,
					Children: make(map[string]domain.Node),
				}
			}
			root.Children[level1].Children[level2].Children[level3].Children[itemID] = domain.Node{
				Item:     true,
				Children: nil,
			}

		} else {
			root.Children[level1].Children[level2].Children[itemID] = domain.Node{
				Item:     true,
				Children: nil,
			}
		}
	}

	return root, nil
}
