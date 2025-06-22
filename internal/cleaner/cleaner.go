package cleaner

import (
	"context"
	"fmt"

	"github.com/DobryySoul/dockr/internal/docker"
	"github.com/docker/docker/api/types/image"
)

func Clean(dckr *docker.DockerClient, images []image.Summary, dryRun bool) error {

	for _, img := range images {
		if dryRun {
			fmt.Printf("[DRY RUN] Would delete image: %s\n", img.ID)
			continue
		}

		_, err := dckr.Cli.ImageRemove(context.Background(), img.ID, image.RemoveOptions{})
		if err != nil {
			return fmt.Errorf("failed to remove image: %w", err)
		}

		fmt.Printf("Deleted image: %s\n", img.ID)
	}
	return nil
}
