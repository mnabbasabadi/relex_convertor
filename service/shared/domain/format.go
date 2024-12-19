// Package domain provides the domain models for the service.
package domain

import (
	"encoding/json"
	"sort"
)

// Data represents a data record.
type Data struct {
	Level1 string
	Level2 string
	Level3 string
	ItemID string
}

// Node represents a node in a tree.
type Node struct {
	Children map[string]Node `json:"children,omitempty"`
	Item     bool            `json:"item,omitempty"`
}

// Marshal marshals the node into a JSON byte slice.
func (n Node) Marshal() ([]byte, error) {
	root := n
	sortChildren := func(children map[string]Node) []string {
		keys := make([]string, 0, len(children))
		for key := range children {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return keys
	}

	type Output struct {
		Children map[string]Node `json:"children"`
	}

	output := Output{
		Children: make(map[string]Node),
	}

	for key, child := range root.Children {
		childKeys := sortChildren(child.Children)
		sortedChildren := make(map[string]Node)
		for _, childKey := range childKeys {
			sortedChildren[childKey] = child.Children[childKey]
		}
		output.Children[key] = Node{
			Item:     false,
			Children: sortedChildren,
		}
	}

	return json.MarshalIndent(output, "", "  ")
}
