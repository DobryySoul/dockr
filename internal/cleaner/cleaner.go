package cleaner

import (
	"context"
	"fmt"

	"github.com/DobryySoul/dockr/internal/docker"
	"github.com/DobryySoul/dockr/internal/domain"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
)

// CleanAll is the main function that triggers the deletion process for all provided unused resources.
// It sequentially calls the cleanup methods for images, containers, volumes, and networks.
func CleanAll(ctx context.Context, client *docker.DockerClient, resources *domain.UnusedResources, all bool) error {
	force := true

	if err := CleanImages(ctx, client, resources.Images); err != nil {
		return err
	}

	if err := CleanContainers(ctx, client, resources.Containers); err != nil {
		return err
	}

	if err := CleanVolumes(ctx, client, resources.Volumes, force); err != nil {
		return err
	}

	if err := CleanNetworks(ctx, client, resources.Networks); err != nil {
		return err
	}

	return nil
}

// CleanImages removes unused (dangling) images.
func CleanImages(ctx context.Context, client *docker.DockerClient, images []*image.Summary) error {
	for _, img := range images {
		_, err := client.Cli.ImageRemove(ctx, img.ID, image.RemoveOptions{})
		if err != nil {
			return fmt.Errorf("failed to remove image with ID: %s, err: %w", img.ID, err)
		}

		fmt.Printf("Deleted image: %s\n", img.ID)
	}

	return nil
}

// CleanContainers removes stopped or dead containers.
func CleanContainers(ctx context.Context, client *docker.DockerClient, containers []*container.Summary) error {
	for _, cont := range containers {

		if err := client.Cli.ContainerRemove(ctx, cont.ID, container.RemoveOptions{}); err != nil {
			return fmt.Errorf("failed to remove container with ID: %s, err: %w", cont.ID, err)
		}

		fmt.Println("Deleted container:", cont.ID)
	}

	return nil
}

// CleanNetworks removes unused networks.
// Ignores system networks and deletes only those not attached to any containers.
func CleanNetworks(ctx context.Context, client *docker.DockerClient, networks []*network.Summary) error {
	for _, net := range networks {
		if err := client.Cli.NetworkRemove(ctx, net.ID); err != nil {
			return fmt.Errorf("failed to remove network with ID: %s, err: %w", net.ID, err)
		}

		fmt.Println("Deleted network:", net.ID)
	}

	return nil
}

// CleanVolumes removes orphaned (unused) data volumes.
// force - forcefully removes the volume (might be needed if Docker still thinks it's busy).
func CleanVolumes(ctx context.Context, client *docker.DockerClient, volumes []*volume.Volume, force bool) error {
	for _, v := range volumes {
		if err := client.Cli.VolumeRemove(ctx, v.Name, force); err != nil {
			return fmt.Errorf("failed to remove volume with name: %s, err: %w", v.Name, err)
		}

		fmt.Println("Deleted volume:", v.Name)
	}

	return nil
}
