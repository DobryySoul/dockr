package docker

import (
	"context"
	"fmt"

	"github.com/DobryySoul/dockr/internal/analyzer"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
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

func (c *DockerClient) FindUnused(ctx context.Context, excludeTags []string) ([]image.Summary, error) {
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

	var unusedImages []image.Summary

	for _, img := range images {
		if analyzer.IsImageUnused(img, excludeTags, usedImages) {
			unusedImages = append(unusedImages, img)
		}
	}

	return unusedImages, nil
}
