package cleaner

import (
	"context"
	"fmt"

	"github.com/DobryySoul/dockr/internal/docker"
	"github.com/DobryySoul/dockr/internal/domain"
	"github.com/docker/docker/api/types/image"
)

func CleanAll(ctx context.Context, client *docker.DockerClient, resources *domain.UnusedResources, all bool) error {
	// Обратите внимание на порядок удаления:
	// 1. Контейнеры
	// 2. Сети (если есть)
	// 3. Образы
	// 4. Тома
	// Это помогает избежать ошибок "ресурс используется".

	if err := CleanImages(ctx, client, resources.Images); err != nil {
		return err
	}

	// if err := CleanContainers(ctx, client, resources.Containers); err != nil {
	// return err
	// }

	// if err := CleanVolumes(ctx, client, resources.Volumes); err != nil {
	// 	return err
	// }

	return nil
}

func CleanImages(ctx context.Context, client *docker.DockerClient, images []*image.Summary) error {
	for _, img := range images {
		_, err := client.Cli.ImageRemove(context.Background(), img.ID, image.RemoveOptions{})
		if err != nil {
			return fmt.Errorf("failed to remove image: %w", err)
		}

		fmt.Printf("Deleted image: %s\n", img.ID)
	}

	return nil
}

// // CleanContainers удаляет контейнеры.
// func CleanContainers(ctx context.Context, client *docker.DockerClient, containers []*container.Summary) error {
// 	// ... логика удаления контейнеров
// }

// // CleanVolumes удаляет тома.
// func CleanVolumes(ctx context.Context, client *docker.DockerClient, volumes []*volume.Volume) error {
// 	// ... логика удаления томов
// }
