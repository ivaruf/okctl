package process_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn/process"
	"github.com/oslokommune/okctl/pkg/mock"
	"github.com/stretchr/testify/assert"
)

const (
	// DefaultSubnetID is a valid, but fake, AWS subnet id
	DefaultSubnetID = "subnet-0bb1c79de3EXAMPLE"
	// DefaultSubnetCIDR is a valid CIDR
	DefaultSubnetCIDR = "192.168.0.0/24"
	// DefaultAvailabilityZone is a valid AWS availability zone
	DefaultAvailabilityZone = "eu-west-1a"
)

func TestSubnets(t *testing.T) {
	testCases := []struct {
		name        string
		provider    v1alpha1.CloudProvider
		value       string
		expect      interface{}
		expectError bool
	}{
		{
			name:     "Should work",
			provider: mock.NewGoodCloudProvider(),
			value:    DefaultSubnetID,
			expect: map[string]v1alpha1.ClusterNetwork{
				DefaultAvailabilityZone: {
					ID:   DefaultSubnetID,
					CIDR: DefaultSubnetCIDR,
				},
			},
		},
		{
			name:        "Should fail",
			provider:    mock.NewBadCloudProvider(),
			value:       DefaultSubnetID,
			expect:      "failed to describe subnet outputs: something bad",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := map[string]v1alpha1.ClusterNetwork{}

			err := process.Subnets(tc.provider, got)(tc.value)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}
