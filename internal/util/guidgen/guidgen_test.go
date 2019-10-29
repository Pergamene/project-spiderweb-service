package guidgen

import (
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/pkg/errors"
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

func TestCheckProposedGUID(t *testing.T) {
	cases := []struct {
		name              string
		paramProposedGUID string
		paramPrefix       string
		paramLength       int
		returnErr         error
	}{
		{
			name:              "happy path",
			paramProposedGUID: "PG_123456789012",
			paramPrefix:       "PG",
			paramLength:       15,
		},
		{
			name:              "happy path just barely enough characters",
			paramProposedGUID: "PG_1",
			paramPrefix:       "PG",
			paramLength:       4,
		},
		{
			name:              "happy path with longer prefix",
			paramProposedGUID: "PGT_12345678901",
			paramPrefix:       "PGT",
			paramLength:       15,
		},
		{
			name:              "invalid length of proposal",
			paramProposedGUID: "PG_12345678901",
			paramPrefix:       "PG",
			paramLength:       15,
			returnErr:         errors.New("proposed guid must be 15 characters"),
		},
		{
			name:              "no prefix",
			paramProposedGUID: "PG_123456789012",
			paramPrefix:       "",
			paramLength:       15,
			returnErr:         errors.New("prefix must be at least one character"),
		},
		{
			name:              "prefix not alphanumeric",
			paramProposedGUID: "",
			paramPrefix:       "PÈ",
			paramLength:       15,
			returnErr:         errors.New("prefix 'PÈ' is not alphanumeric. Do not include '_' in the prefix. It will be added automatically"),
		},
		{
			name:              "prefix not alphanumeric due to added _",
			paramProposedGUID: "",
			paramPrefix:       "PG_",
			paramLength:       15,
			returnErr:         errors.New("prefix 'PG_' is not alphanumeric. Do not include '_' in the prefix. It will be added automatically"),
		},
		{
			name:              "length too short",
			paramProposedGUID: "12",
			paramPrefix:       "PG",
			paramLength:       2,
			returnErr:         errors.New("length of 2 is invalid. Must be at least 3 for a single-character prefix + '_' + a random character for the guid"),
		},
		{
			name:              "length too short due to prefix size",
			paramProposedGUID: "",
			paramPrefix:       "PGXYZ",
			paramLength:       6,
			returnErr:         errors.New("proposed prefix must be less than 4 characters"),
		},
		{
			name:              "proposal does not have proper prefix",
			paramProposedGUID: "PGT_12345678901",
			paramPrefix:       "PG",
			paramLength:       15,
			returnErr:         errors.New("proposed guid must start with 'PG_'"),
		},
		{
			name:              "proposal is not alphanumeric",
			paramProposedGUID: "PG_12345678901È",
			paramPrefix:       "PG",
			paramLength:       15,
			returnErr:         errors.New("characters after prefix, 'PG_', must be English alphanumeric"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckProposedGUID(tc.paramProposedGUID, tc.paramPrefix, tc.paramLength)
			testutils.TestErrorAgainstCase(t, err, tc.returnErr)
		})
	}
}
