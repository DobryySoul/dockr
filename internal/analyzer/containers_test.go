package analyzer

import (
	"testing"

	"github.com/docker/docker/api/types/container"
)

func TestIsContainerUnused(t *testing.T) {
	tests := []struct {
		name     string
		state    string
		expected bool
	}{
		{
			name:     "running container",
			state:    "running",
			expected: false,
		},
		{
			name:     "exited container",
			state:    "exited",
			expected: true,
		},
		{
			name:     "created container",
			state:    "created",
			expected: true,
		},
		{
			name:     "dead container",
			state:    "dead",
			expected: true,
		},
		{
			name:     "paused container",
			state:    "paused",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &container.Summary{
				State: tt.state,
			}
			result := IsContainerUnused(c)
			if result != tt.expected {
				t.Errorf("expected %v for state %s, got %v", tt.expected, tt.state, result)
			}
		})
	}
}
