package analyzer

import (
	"testing"

	"github.com/docker/docker/api/types/image"
)

func TestIsImageUnused(t *testing.T) {
	tests := []struct {
		name        string
		img         image.Summary
		excludeTags []string
		usedImages  map[string]bool
		expected    bool
	}{
		{
			name: "used image",
			img: image.Summary{
				ID:       "image-id-1",
				RepoTags: []string{"ubuntu:latest"},
			},
			excludeTags: []string{},
			usedImages: map[string]bool{
				"image-id-1": true,
			},
			expected: false,
		},
		{
			name: "unused image with no tags",
			img: image.Summary{
				ID:       "image-id-2",
				RepoTags: []string{},
			},
			excludeTags: []string{},
			usedImages:  map[string]bool{},
			expected:    true,
		},
		{
			name: "unused image but in exclude list",
			img: image.Summary{
				ID:       "image-id-3",
				RepoTags: []string{"my-app:prod"},
			},
			excludeTags: []string{"prod"},
			usedImages:  map[string]bool{},
			expected:    false,
		},
		{
			name: "unused image not in exclude list",
			img: image.Summary{
				ID:       "image-id-4",
				RepoTags: []string{"my-app:dev"},
			},
			excludeTags: []string{"prod"},
			usedImages:  map[string]bool{},
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsImageUnused(tt.img, tt.excludeTags, tt.usedImages)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
