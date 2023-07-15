package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode_Marshal(t *testing.T) {

	tests := []struct {
		name string
		node Node
		want []byte
	}{
		{
			name: "Valid payload",
			node: Node{
				Children: map[string]Node{
					"A": {Children: map[string]Node{
						"B": {Children: map[string]Node{
							"C": {Children: map[string]Node{
								"1": {
									Item: true,
								},
								"2": {
									Item: true,
								},
							}},
							"D": {Children: map[string]Node{
								"3": {
									Item: true,
								},
							}},
						}},
						"E": {Children: map[string]Node{
							"F": {Children: map[string]Node{
								"4": {
									Item: true,
								},
							}},
						}},
					}},
				},
			},
			want: []byte(`{
  "children": {
    "A": {
      "children": {
        "B": {
          "children": {
            "C": {
              "children": {
                "1": {
                  "item": true
                },
                "2": {
                  "item": true
                }
              }
            },
            "D": {
              "children": {
                "3": {
                  "item": true
                }
              }
            }
          }
        },
        "E": {
          "children": {
            "F": {
              "children": {
                "4": {
                  "item": true
                }
              }
            }
          }
        }
      }
    }
  }
}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.node.Marshal()
			require.NoError(t, err)
			require.Equal(t, tt.want, got, "Marshal() = %v, want %v", string(got), string(tt.want))
		})
	}
}
