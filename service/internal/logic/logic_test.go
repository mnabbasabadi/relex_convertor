package logic

import (
	"context"
	"os"
	"testing"

	"github.com/mnabbasabadi/relex_convertor/service/shared/domain"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestConvertToJSON(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name     string
		input    []domain.Data
		in       [][]string
		expected domain.Node
		wantErr  bool
	}{
		{
			name: "Valid payload",
			in: [][]string{
				{"level_1", "level_2", "level_3", "item_id"},
				{"A", "B", "C", "1"},
				{"A", "B", "C", "2"},
				{"A", "B", "D", "3"},
				{"A", "E", "F", "4"},
			},
			expected: domain.Node{
				Children: map[string]domain.Node{
					"A": {Children: map[string]domain.Node{
						"B": {Children: map[string]domain.Node{
							"C": {Children: map[string]domain.Node{
								"1": {
									Item: true,
								},
								"2": {
									Item: true,
								},
							}},
							"D": {Children: map[string]domain.Node{
								"3": {
									Item: true,
								},
							}},
						}},
						"E": {Children: map[string]domain.Node{
							"F": {Children: map[string]domain.Node{
								"4": {
									Item: true,
								},
							}},
						}},
					}},
				},
			},
		},
		{
			name: "Invalid payload: empty level_1 with non-empty level_2",
			in: [][]string{
				{"level_1", "level_2", "level_3", "item_id"},
				{"", "B", "C", "1"},
				{"A", "", "C", "2"},
			},
			wantErr: true,
		},
		{
			name: "Invalid payload: empty level_2 with non-empty item_id",
			in: [][]string{
				{"level_1", "level_2", "level_3", "item_id"},
				{"A", "", "C", "1"},
				{"A", "B", "", "2"},
			},
			wantErr: true,
		},
	}

	// Create a mock logger
	logger := zerolog.New(os.Stdout)

	controller := New(logger)

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := controller.ConvertToTreeNode(context.Background(), tc.in)
			require.Equal(t, tc.wantErr, err != nil)
			if tc.wantErr {
				return
			}
			require.Equal(t, tc.expected, result)
		})
	}
}
