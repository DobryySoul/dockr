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

// CleanAll - основная функция, запускающая процесс удаления всех переданных неиспользуемых ресурсов.
// Поочередно вызывает методы очистки для образов, контейнеров, томов и сетей.
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

// CleanImages удаляет неиспользуемые ("dangling") образы.
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

// CleanContainers удаляет остановленные (или мертвые) контейнеры.
func CleanContainers(ctx context.Context, client *docker.DockerClient, containers []*container.Summary) error {
	for _, cont := range containers {

		if err := client.Cli.ContainerRemove(ctx, cont.ID, container.RemoveOptions{}); err != nil {
			return fmt.Errorf("failed to remove container with ID: %s, err: %w", cont.ID, err)
		}

		fmt.Println("Deleted container:", cont.ID)
	}

	return nil
}

// CleanNetworks удаляет неиспользуемые сети.
// Игнорирует системные сети и удаляет только те, которые не привязаны к контейнерам.
func CleanNetworks(ctx context.Context, client *docker.DockerClient, networks []*network.Summary) error {
	for _, net := range networks {
		if err := client.Cli.NetworkRemove(ctx, net.ID); err != nil {
			return fmt.Errorf("failed to remove network with ID: %s, err: %w", net.ID, err)
		}

		fmt.Println("Deleted network:", net.ID)
	}

	return nil
}

// CleanVolumes удаляет "осиротевшие" (неиспользуемые) тома данных.
// force - принудительное удаление (может потребоваться, если Docker все еще думает, что том занят).
func CleanVolumes(ctx context.Context, client *docker.DockerClient, volumes []*volume.Volume, force bool) error {
	for _, v := range volumes {
		if err := client.Cli.VolumeRemove(ctx, v.Name, force); err != nil {
			return fmt.Errorf("failed to remove volume with name: %s, err: %w", v.Name, err)
		}

		fmt.Println("Deleted volume:", v.Name)
	}

	return nil
}
