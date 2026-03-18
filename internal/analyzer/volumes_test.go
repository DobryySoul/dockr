package analyzer

import (
	"testing"

	"github.com/docker/docker/api/types/volume"
)

func TestIsVolumeUnused(t *testing.T) {
	tests := []struct {
		name        string
		vol         *volume.Volume
		usedVolumes map[string]bool
		expected    bool
	}{
		{
			name: "used volume",
			vol: &volume.Volume{
				Name: "my-volume-1",
			},
			usedVolumes: map[string]bool{
				"my-volume-1": true,
			},
			expected: false,
		},
		{
			name: "unused volume",
			vol: &volume.Volume{
				Name: "my-volume-2",
			},
			usedVolumes: map[string]bool{
				"my-volume-1": true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsVolumeUnused(tt.vol, tt.usedVolumes)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
