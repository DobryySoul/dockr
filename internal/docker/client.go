package docker

import (
	"context"
	"fmt"

	"github.com/DobryySoul/dockr/internal/analyzer"
	"github.com/DobryySoul/dockr/internal/domain"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	Cli *client.Client
}

// NewDockerClient создает новый клиент для взаимодействия с Docker API.
// Клиент автоматически настраивает версию API для совместимости с хостом.
func NewDockerClient(ctx context.Context) (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &DockerClient{Cli: cli}, nil
}

// FindUnusedResourcer собирает все неиспользуемые Docker-ресурсы (образы, контейнеры, тома, сети),
// которые можно безопасно удалить. Возвращает структуру domain.UnusedResources.
func (c *DockerClient) FindUnusedResourcer(ctx context.Context, excludeTags []string) (*domain.UnusedResources, error) {
	images, err := c.FindUnusedImages(ctx, excludeTags)
	if err != nil {
		return nil, err
	}

	containers, err := c.FindUnusedContainers(ctx)
	if err != nil {
		return nil, err
	}

	volumes, err := c.FindUnusedVolumes(ctx, false)
	if err != nil {
		return nil, err
	}

	networks, err := c.FindUnusedNetworks(ctx)
	if err != nil {
		return nil, err
	}

	resources := &domain.UnusedResources{
		Images:     images,
		Containers: containers,
		Volumes:    volumes,
		Networks:   networks,
	}

	return resources, nil
}

// FindUnusedImages ищет неиспользуемые (dangling) образы. Образ считается неиспользуемым,
// если к нему не привязан ни один контейнер и его тег не находится в списке excludeTags.
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

// FindUnusedContainers ищет остановленные, созданные или "мертвые" контейнеры, 
// которые больше не выполняют никакой работы и занимают место на диске.
func (c *DockerClient) FindUnusedContainers(ctx context.Context) ([]*container.Summary, error) {
	containers, err := c.Cli.ContainerList(ctx, container.ListOptions{All: true, Size: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list Docker containers: %w", err)
	}

	var unusedContainers []*container.Summary
	for _, cont := range containers {
		contCopy := cont
		if analyzer.IsContainerUnused(&contCopy) {
			unusedContainers = append(unusedContainers, &contCopy)
		}
	}

	return unusedContainers, nil
}

// FindUnusedNetworks ищет неиспользуемые сети.
func (c *DockerClient) FindUnusedNetworks(ctx context.Context) ([]*network.Summary, error) {
	networks, err := c.Cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Docker networks: %w", err)
	}

	var unusedNetworks []*network.Summary
	for _, net := range networks {
		netCopy := net
		if analyzer.IsNetworkUnused(&netCopy) {
			unusedNetworks = append(unusedNetworks, &netCopy)
		}
	}

	return unusedNetworks, nil
}

// FindUnusedVolumes ищет "осиротевшие" (dangling) тома.
// Том считается неиспользуемым, если он не примонтирован ни к одному из существующих контейнеров.
func (c *DockerClient) FindUnusedVolumes(ctx context.Context, force bool) ([]*volume.Volume, error) {
	containers, err := c.Cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list Docker containers: %w", err)
	}

	usedVolumes := make(map[string]bool)
	for _, c := range containers {
		for _, mount := range c.Mounts {
			if mount.Type == "volume" {
				usedVolumes[mount.Name] = true
			}
		}
	}

	volumesList, err := c.Cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Docker volumes: %w", err)
	}

	var unusedVolumes []*volume.Volume
	for _, v := range volumesList.Volumes {
		vCopy := v
		if analyzer.IsVolumeUnused(vCopy, usedVolumes) {
			unusedVolumes = append(unusedVolumes, vCopy)
		}
	}

	return unusedVolumes, nil
}
