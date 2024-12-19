package converter

import (
	"context"
	"testing"

	"github.com/mnabbasabadi/relex_convertor/service/shared/domain"
	"github.com/stretchr/testify/require"
)

func TestConvert_BuildTree(t *testing.T) {
	testCases := map[string]struct {
		in       []domain.Data
		expected domain.Node
	}{
		"Valid payload": {
			in: []domain.Data{
				{
					Level1: "A",
					Level2: "B",
					Level3: "C",
					ItemID: "1",
				},
				{
					Level1: "A",
					Level2: "B",
					Level3: "C",
					ItemID: "2",
				},
				{
					Level1: "A",
					Level2: "B",
					Level3: "D",
					ItemID: "3",
				},
				{
					Level1: "A",
					Level2: "E",
					Level3: "F",
					ItemID: "4",
				},
			},
			expected: domain.Node{
				Children: map[string]domain.Node{
					"A": {
						Children: map[string]domain.Node{
							"B": {
								Children: map[string]domain.Node{
									"C": {
										Children: map[string]domain.Node{
											"1": {Item: true},
											"2": {Item: true},
										},
									},
									"D": {
										Children: map[string]domain.Node{
											"3": {Item: true},
										},
									},
								},
							},
							"E": {
								Children: map[string]domain.Node{
									"F": {
										Children: map[string]domain.Node{
											"4": {Item: true},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"valid payload with empty level_2": {
			in: []domain.Data{
				{
					Level1: "A",
					ItemID: "1",
				},
				{
					Level1: "A",
					ItemID: "2",
				},
			},
			expected: domain.Node{
				Children: map[string]domain.Node{
					"A": {
						Children: map[string]domain.Node{
							"1": {Item: true},
							"2": {Item: true},
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := New().BuildTree(context.Background(), tc.in)
			require.NoError(t, err)
			require.Equal(t, tc.expected, got)

		})
	}
}
