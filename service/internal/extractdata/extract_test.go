package extractdata

import (
	"context"
	"testing"

	"github.com/mnabbasabadi/relex_convertor/service/shared/domain"
	"github.com/stretchr/testify/require"
)

func TestController_Extract(t *testing.T) {
	testCases := map[string]struct {
		input   [][]string
		want    []domain.Data
		wantErr bool
	}{
		"Valid payload": {
			input: [][]string{
				{"level_1", "level_2", "level_3", "item_id"},
				{"A", "B", "C", "1"},
				{"A", "B", "C", "2"},
				{"A", "B", "D", "3"},
			},
			want: []domain.Data{
				{Level1: "A", Level2: "B", Level3: "C", ItemID: "1"},
				{Level1: "A", Level2: "B", Level3: "C", ItemID: "2"},
				{Level1: "A", Level2: "B", Level3: "D", ItemID: "3"},
			},
		},
		"Invalid payload: empty level_1 with non-empty level_2": {
			input: [][]string{
				{"level_1", "level_2", "level_3", "item_id"},
				{"", "B", "C", "1"},
				{"A", "", "C", "2"},
			},
			wantErr: true,
		},
		"Invalid payload: empty level_2 with non-empty level_3": {
			input: [][]string{
				{"level_1", "level_2", "level_3", "item_id"},
				{"A", "", "C", "1"},
				{"A", "B", "", "2"},
			},
			wantErr: true,
		},
		"invalid payload: number of records is less than 2": {
			input: [][]string{
				{"level_1"},
				{"A"},
			},
			wantErr: true,
		},
		"empty payload": {
			input:   [][]string{},
			wantErr: true,
		},
		"one record payload": {
			input: [][]string{
				{"level_1", "item_id"},
				{"A", "1"},
			},
			want: []domain.Data{
				{Level1: "A", ItemID: "1"},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := New().Extract(context.Background(), tc.input)
			require.Equal(t, tc.wantErr, err != nil)
			if tc.wantErr {
				return
			}
			require.Equal(t, tc.want, got)
		})
	}
}
