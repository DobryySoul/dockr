package analyzer

import "github.com/docker/docker/api/types/volume"

// IsVolumeUnused проверяет, является ли том неиспользуемым (осиротевшим).
// В Docker API тома, не прикрепленные к контейнерам, считаются неиспользуемыми.
func IsVolumeUnused(vol *volume.Volume, usedVolumes map[string]bool) bool {
	return !usedVolumes[vol.Name]
}
