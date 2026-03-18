package analyzer

import (
	"testing"

	"github.com/docker/docker/api/types/network"
)

func TestIsNetworkUnused(t *testing.T) {
	tests := []struct {
		name     string
		net      *network.Summary
		expected bool
	}{
		{
			name: "bridge network",
			net: &network.Summary{
				Name: "bridge",
			},
			expected: false,
		},
		{
			name: "host network",
			net: &network.Summary{
				Name: "host",
			},
			expected: false,
		},
		{
			name: "none network",
			net: &network.Summary{
				Name: "none",
			},
			expected: false,
		},
		{
			name: "custom network with containers",
			net: &network.Summary{
				Name: "my-net",
				Containers: map[string]network.EndpointResource{
					"container1": {},
				},
			},
			expected: false,
		},
		{
			name: "custom network without containers",
			net: &network.Summary{
				Name:       "my-net-unused",
				Containers: map[string]network.EndpointResource{},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNetworkUnused(tt.net)
			if result != tt.expected {
				t.Errorf("expected %v for network %s, got %v", tt.expected, tt.net.Name, result)
			}
		})
	}
}
