package guidgen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateGUID(t *testing.T) {
	cases := []struct {
		name               string
		paramPrefix        string
		paramLength        int
		returnStringPrefix string
		returnStringLength int
	}{
		{
			name:               "test generation of typical guid",
			paramPrefix:        "PG",
			paramLength:        15,
			returnStringPrefix: "PG_",
			returnStringLength: 15,
		},
		{
			name:               "test generation of too small of guid",
			paramPrefix:        "ABCDE",
			paramLength:        3,
			returnStringPrefix: "ABCDE",
			returnStringLength: len("ABCDE"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := GenerateGUID(tc.paramPrefix, tc.paramLength)
			require.Equal(t, tc.returnStringLength, len(result))
			require.Equal(t, tc.returnStringPrefix, result[:len(tc.returnStringPrefix)])
		})
	}
}
