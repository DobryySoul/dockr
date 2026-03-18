package analyzer

import "github.com/docker/docker/api/types/volume"

// IsVolumeUnused checks if a volume is unused (orphaned).
// Volumes that are not attached to any container are considered unused.
func IsVolumeUnused(vol *volume.Volume, usedVolumes map[string]bool) bool {
	return !usedVolumes[vol.Name]
}
