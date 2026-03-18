package analyzer

import "github.com/docker/docker/api/types/container"

// IsContainerUnused checks if the container is unused (stopped).
// Containers with "exited", "created", or "dead" states can be safely removed.
func IsContainerUnused(c *container.Summary) bool {
	// Valid states for deletion
	return c.State == "exited" || c.State == "created" || c.State == "dead"
}
