package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/DobryySoul/dockr/internal/cleaner"
	"github.com/DobryySoul/dockr/internal/docker"
	"github.com/DobryySoul/dockr/internal/formatter"
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
	Short: "Smart Docker resource cleaner",
	Long: `Utility for safely removing unused Docker resources:
- Images
- Containers
- Volumes
- Networks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if version {
			fmt.Println(versionApp)
		}

		dockerClient, err := docker.NewDockerClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to connect to Docker: %w", err)
		}

		resources, err := dockerClient.FindUnusedResourcer(ctx, excludeTags)
		if err != nil {
			return fmt.Errorf("analysis error: %w", err)
		}

		if resources.IsEmpty() {
			formatter.Info("No unused resources found.")
			return nil
		}

		formatter.PrintReport(resources, dryRun)

		if dryRun {
			return nil
		}

		if interactive && !formatter.Confirm("Proceed with deletion?", resources) {
			formatter.Info("Operation cancelled")
		}

		if err := cleaner.CleanAll(ctx, dockerClient, resources, all); err != nil {
			return fmt.Errorf("cleanup error: %w", err)
		}

		formatter.Success("Cleanup completed! Reclaimed: %.2f MB",
			float64(resources.TotalSize()/mb))

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Show version")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Simulate deletion without actually removing resources")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Prompt for confirmation before removing resources")
	rootCmd.Flags().StringSliceVarP(&excludeTags, "exclude-tags", "e", []string{}, "List of image tags to exclude from deletion")
	rootCmd.Flags().BoolVarP(&all, "all", "a", false, "Remove ALL unused resources (including important ones)")
}
