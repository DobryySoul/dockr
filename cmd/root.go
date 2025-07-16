package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/DobryySoul/dockr/internal/cleaner"
	"github.com/DobryySoul/dockr/internal/docker"
	"github.com/DobryySoul/dockr/pkg/formatter"
	"github.com/spf13/cobra"
)

const (
	mb         = 1024 * 1024
	versionApp = "dockr version v1.0.0"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if version {
			fmt.Println(versionApp)
		}

		dockerClient, err := docker.NewDockerClient(ctx)
		if err != nil {
			return fmt.Errorf("ошибка подключения к Docker: %w", err)
		}

		resources, err := dockerClient.FindUnusedResourcer(ctx, excludeTags)
		if err != nil {
			return fmt.Errorf("ошибка анализа: %w", err)
		}

		if resources.IsEmpty() {
			formatter.Info("Неиспользуемых ресурсов не найдено.")
			return nil
		}

		formatter.PrintReport(resources, dryRun)

		if dryRun {
			return nil
		}

		if interactive && !formatter.Confirm("Продолжить удаление?", resources) {
			formatter.Info("Операция отменена")
		}

		if err := cleaner.CleanAll(ctx, dockerClient, resources, all); err != nil {
			return fmt.Errorf("ошибка очистки: %w", err)
		}

		formatter.Success("Очистка завершена! Освобождено: %.2f MB",
			float64(resources.TotalSize()/mb))

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Просмотр версии")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Выводит информацию о ресурсах, которые будут удалены")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Запрашивает подтверждение перед удалением ресурсов")
	rootCmd.Flags().StringSliceVarP(&excludeTags, "exclude-tags", "e", []string{}, "Список тегов образов, которые нужно исключить при удалении")
	rootCmd.Flags().BoolVarP(&all, "all", "a", false, "Удалять ВСЕ неиспользуемые ресурсы (включая важные)")
}
