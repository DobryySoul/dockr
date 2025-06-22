package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/DobryySoul/dockr/internal/cleaner"
	"github.com/DobryySoul/dockr/internal/docker"
	"github.com/DobryySoul/dockr/internal/domain"
	"github.com/spf13/cobra"
)

var (
	dryRun      bool
	interactive bool
	excludeTags []string
	all         bool
)

var rootCmd = &cobra.Command{
	Use:   "docker-cleaner",
	Short: "Умный очиститель Docker-ресурсов",
	Long: `Утилита для безопасного удаления неиспользуемых Docker-ресурсов:
- Образы (images)
- Контейнеры (containers)
- Тома (volumes)
- Сети (networks)`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dockerClient, err := docker.NewDockerClient(ctx)
		if err != nil {
			fmt.Printf("Ошибка подключения к Docker: %v\n", err)
			return
		}

		images, err := dockerClient.FindUnusedImages(ctx, excludeTags)
		if err != nil {
			fmt.Printf("Ошибка при получении результатов анализа: %v\n", err)
			return
		}

		resourcer := &domain.UnusedResources{Images: images}

		if dryRun {
			printDryRunResults(resourcer)
			return
		}

		if interactive {
			if !confirmDeletion(resourcer) {
				fmt.Println("Отмена операции")
				return
			}
		}

		if err := cleaner.RemoveImages(dockerClient, resourcer.Images, all); err != nil {
			fmt.Printf("Ошибка очистки: %v\n", err)
		} else {
			fmt.Println("Очистка завершена успешно")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Выводит информацию о ресурсах, которые будут удалены")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Запрашивает подтверждение перед удалением ресурсов")
	rootCmd.Flags().StringSliceVarP(&excludeTags, "exclude-tags", "e", []string{}, "Список тегов образов, которые нужно исключить при удалении")
	rootCmd.Flags().BoolVarP(&all, "all", "a", false, "Удалять ВСЕ неиспользуемые ресурсы (включая важные)")
}

func printDryRunResults(resources *domain.UnusedResources) {
	fmt.Println("\n=== Ресурсы для удаления (dry-run) ===")

	if len(resources.Images) > 0 {
		fmt.Println("\nОбразы:")
		for _, img := range resources.Images {
			tags := "<none>"
			if len(img.RepoTags) > 0 {
				tags = strings.Join(img.RepoTags, ", ")
			}
			fmt.Printf("- %-12s %-40s %8.2f MB\n",
				img.ID[:12],
				truncate(tags, 40),
				float64(img.Size)/1024/1024)
		}
	}

	// if len(resources.Containers) > 0 {
	// 	fmt.Println("\nКонтейнеры:")
	// 	for _, c := range resources.Containers {
	// 		fmt.Printf("- %-12s %-20s %s\n",
	// 			c.ID[:12],
	// 			truncate(c.Names[0], 20),
	// 			c.State)
	// 	}
	// }

	fmt.Printf("\nИтого будет освобождено: %.2f MB\n",
		resources.TotalSize()/1024/1024)
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

func confirmDeletion(resources *domain.UnusedResources) bool {
	fmt.Printf("\nБудет удалено:\n")
	fmt.Printf("- Образы: %d (%.2f MB)\n", len(resources.Images), resources.ImagesSize()/1024/1024)
	// fmt.Printf("- Контейнеры: %d\n", len(resources.Containers))
	// fmt.Printf("- Тома: %d\n", len(resources.Volumes))
	// fmt.Printf("- Сети: %d\n", len(resources.Networks))

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nПродолжить удаление? [y/N]: ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))

		switch answer {
		case "y", "yes":
			return true
		case "n", "no", "":
			return false
		default:
			fmt.Println("Пожалуйста, введите 'y' или 'n'")
		}
	}
}
