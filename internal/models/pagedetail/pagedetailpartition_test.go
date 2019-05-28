package pagedetail

import (
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/require"
)

func TestUnmarshelPartitions(t *testing.T) {
	cases := []struct {
		name                string
		paramPartitions     []Partition
		resultingPartitions []Partition
		returnErr           error
	}{
		{
			name: "verify recursive partitions with items",
			paramPartitions: []Partition{
				{
					TypeString: "h1",
					Partitions: []Partition{
						{
							TypeString: "h2",
							Items: []Partition{
								{
									TypeString: "h3",
								},
								{
									TypeString: "h4",
								},
							},
						},
					},
				},
			},
			resultingPartitions: []Partition{
				{
					TypeString: "h1",
					Type:       PartitionTypeHeaderOne,
					Partitions: []Partition{
						{
							TypeString: "h2",
							Type:       PartitionTypeHeaderTwo,
							Items: []Partition{
								{
									TypeString: "h3",
									Type:       PartitionTypeHeaderThree,
								},
								{
									TypeString: "h4",
									Type:       PartitionTypeHeaderFour,
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := UnmarshelPartitions(tc.paramPartitions)
			testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if err != nil {
				return
			}
			require.Equal(t, tc.resultingPartitions, tc.paramPartitions)
		})
	}
}
