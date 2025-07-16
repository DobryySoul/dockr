package docker

import (
	"context"
	"fmt"

	"github.com/DobryySoul/dockr/internal/analyzer"
	"github.com/DobryySoul/dockr/internal/domain"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	Cli *client.Client
}

func NewDockerClient(ctx context.Context) (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &DockerClient{Cli: cli}, nil
}

func (c *DockerClient) FindUnusedResourcer(ctx context.Context, excludeTags []string) (*domain.UnusedResources, error) {
	images, err := c.FindUnusedImages(ctx, excludeTags)
	if err != nil {
		return nil, err
	}

	// containers, err := c.FindUnusedContainers(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	volumes, err := c.FindUnusedVolumes(ctx, false)
	if err != nil {
		return nil, err
	}

	// Поиск сетей - аналогично
	// networks, err := c.FindUnusedNetworks(ctx) ...

	resources := &domain.UnusedResources{
		Images: images,
		// Containers: containers,
		Volumes: volumes,
		// Networks: networks,
	}

	return resources, nil
}

func (c *DockerClient) FindUnusedImages(ctx context.Context, excludeTags []string) ([]*image.Summary, error) {
	containers, err := c.Cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list Docker containers: %w", err)
	}

	images, err := c.Cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Docker images: %w", err)
	}

	usedImages := make(map[string]bool)
	for _, container := range containers {
		usedImages[container.ImageID] = true
	}

	var unusedImages []*image.Summary

	for _, img := range images {
		if analyzer.IsImageUnused(img, excludeTags, usedImages) {
			unusedImages = append(unusedImages, &img)
		}
	}

	return unusedImages, nil
}

// FindUnusedContainers ищет остановленные контейнеры.
// func (c *DockerClient) FindUnusedContainers(ctx context.Context) ([]*container.Summary, error) {
// ... логика для поиска контейнеров
// (docker ps -a -f "status=exited" -f "status=created")
// возвращает []types.Container, error
// }

// FindUnusedVolumes ищет "осиротевшие" тома.
func (c *DockerClient) FindUnusedVolumes(ctx context.Context, force bool) ([]*volume.Volume, error) {
	volumesList, err := c.Cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Docker volumes: %w", err)
	}

	return volumesList.Volumes, nil
}
