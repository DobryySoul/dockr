package analyzer

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/image"
)

func IsImageUnused(img image.Summary, excludeTags []string, usedImages map[string]bool) bool {
	if usedImages[img.ID] {
		return false
	}

	if len(img.RepoTags) == 0 {
		return true
	}

	for _, tag := range img.RepoTags {
		for _, excluded := range excludeTags {
			fmt.Printf("image tag [%s] - > исключение [%s]", tag, excluded)
			if strings.Contains(tag, excluded) {
				return false
			}
		}
	}

	return true
}
