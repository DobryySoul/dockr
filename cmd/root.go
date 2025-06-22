package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/DobryySoul/dockr/internal/cleaner"
	"github.com/DobryySoul/dockr/internal/docker"
	"github.com/DobryySoul/dockr/internal/domain"
	"github.com/DobryySoul/dockr/pkg/formatter"
	"github.com/spf13/cobra"
)

var (
	dryRun      bool
	interactive bool
	excludeTags []string
	all         bool
	version     bool
)

var rootCmd = &cobra.Command{
	Use:   "dockr",
	Short: "Умный очиститель Docker-ресурсов",
	Long: `Утилита для безопасного удаления неиспользуемых Docker-ресурсов:
- Образы (images)
- Контейнеры (containers)
- Тома (volumes)
- Сети (networks)`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if version {
			fmt.Println("dockr v1.0.0")
			return
		}

		dockerClient, err := docker.NewDockerClient(ctx)
		if err != nil {
			formatter.Error("Ошибка подключения к Docker: %v", err)
			return
		}

		images, err := dockerClient.FindUnused(ctx, excludeTags)
		if err != nil {
			formatter.Error("Ошибка анализа: %v", err)
			return
		}

		resources := &domain.UnusedResources{Images: images}

		formatter.PrintReport(resources, dryRun)

		if dryRun {
			return
		}

		if interactive && !formatter.Confirm("Продолжить удаление?", resources) {
			formatter.Info("Операция отменена")
			return
		}

		if err := cleaner.Clean(dockerClient, resources.Images, all); err != nil {
			formatter.Error("Ошибка очистки: %v", err)
			os.Exit(1)
		}

		formatter.Success("Очистка завершена! Освобождено: %.2f MB",
			(*domain.UnusedResources).TotalSize(resources)/1024/1024)
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
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Просмотр версии")
}
