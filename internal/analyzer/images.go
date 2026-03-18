package analyzer

import (
	"strings"

	"github.com/docker/docker/api/types/image"
)

// IsImageUnused checks if an image is considered dangling/unused based on its tags
// and whether it is currently used by any container.
func IsImageUnused(img image.Summary, excludeTags []string, usedImages map[string]bool) bool {
	if usedImages[img.ID] {
		return false
	}

	if len(img.RepoTags) == 0 {
		return true
	}

	for _, tag := range img.RepoTags {
		for _, excluded := range excludeTags {
			if strings.Contains(tag, excluded) {
				return false
			}
		}
	}

	return true
}
