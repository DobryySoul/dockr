//nolint:errcheck,gosec // We intentionally ignore error returns from formatting/printing functions
package formatter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/DobryySoul/dockr/internal/domain"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"

	"github.com/fatih/color"
)

var (
	ErrorColor   = color.New(color.FgRed)
	SuccessColor = color.New(color.FgGreen)
	InfoColor    = color.New(color.FgBlue)
	WarningColor = color.New(color.FgYellow)
)

func Error(format string, a ...any) {
	ErrorColor.Printf("✖ "+format+"\n", a...)
}

func Success(format string, a ...any) {
	SuccessColor.Printf("✔ "+format+"\n", a...)
}

func Info(format string, a ...any) {
	InfoColor.Printf("ℹ "+format+"\n", a...)
}

// Confirm displays a summary of how much space will be freed and prompts
// the user for confirmation to proceed (y/n). Returns true if the user agreed.
func Confirm(question string, resources *domain.UnusedResources) bool {
	fmt.Printf("\nWill be deleted:\n")
	fmt.Printf("- Images: %d (%.2f MB)\n", len(resources.Images), resources.ImagesSize()/1024/1024)
	fmt.Printf("- Containers: %d (%.2f MB)\n", len(resources.Containers), resources.ContainersSize()/1024/1024)
	fmt.Printf("- Volumes: %d (%.2f MB)\n", len(resources.Volumes), resources.VolumesSize()/1024/1024)
	fmt.Printf("- Networks: %d\n", len(resources.Networks))

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nProceed with deletion? [y/N]: ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))

		switch answer {
		case "y", "yes":
			return true
		case "n", "no", "":
			return false
		default:
			fmt.Println("Please type 'y' or 'n'")
		}
	}
}

// PrintReport prints a detailed structured report (in tables)
// of all found unused resources. In dryRun mode, it prints a warning header.
func PrintReport(res *domain.UnusedResources, dryRun bool) {
	if dryRun {
		color.New(color.FgYellow).Println("\n=== DRY RUN MODE ===")
		color.New(color.FgHiBlue).Println("The following resources will be deleted:")
	} else {
		color.New(color.FgYellow).Println("\n=== DOCKER CLEANUP ===")
	}

	printSection := func(title string, count int, fn func()) {
		if count > 0 {
			color.New(color.FgGreen).Printf("\n%s (%d):\n", title, count)
			fn()
		}
	}

	printSection("Images", len(res.Images), func() {
		printImagesTable(res.Images)
	})

	printSection("Containers", len(res.Containers), func() {
		printContainersTable(res.Containers)
	})

	printSection("Volumes", len(res.Volumes), func() {
		printVolumesTable(res.Volumes)
	})

	printSection("Networks", len(res.Networks), func() {
		printNetworksTable(res.Networks)
	})

	if res.TotalCount() > 0 {
		color.New(color.FgHiWhite).Printf("\nTotal: %d resources, ", res.TotalCount())
		color.New(color.FgHiGreen).Printf("freed space: %.2f MB\n", res.TotalSize()/1024/1024)
	} else {
		color.New(color.FgHiGreen).Println("\nNothing to delete - system is clean!")
	}
}

func printImagesTable(images []*image.Summary) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\t TAG\t SIZE\t")

	for _, img := range images {
		tags := strings.Join(img.RepoTags, ", ")
		if tags == "" {
			tags = color.HiRedString("<none>")
		}

		fmt.Fprintf(w, "%s\t %s\t %.2f MB\n",
			truncateID(img.ID),
			truncate(tags, 30),
			float64(img.Size)/1024/1024,
		)
	}
	w.Flush()
}

func printContainersTable(containers []*container.Summary) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\t NAME\t STATE\t IMAGE\t SIZE\t")

	for _, c := range containers {
		state := c.State
		if c.State == "exited" || c.State == "dead" {
			state = color.HiRedString(state)
		}

		fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %.2f MB\t\n",
			truncateID(c.ID),
			truncate(c.Names[0], 20),
			state,
			truncate(c.Image, 20),
			float64(c.SizeRw)/1024/1024,
		)
	}
	w.Flush()
}

func printVolumesTable(volumes []*volume.Volume) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DRIVER\t NAME\t SIZE\t")

	for _, v := range volumes {
		size := 0.0
		if v.UsageData != nil {
			size = float64(v.UsageData.Size) / 1024 / 1024
		}

		fmt.Fprintf(w, "%s\t %s\t %.2f MB\t\n",
			truncateID(v.Driver),
			truncate(v.Name, 40),
			size,
		)
	}
	w.Flush()
}

func printNetworksTable(networks []*network.Summary) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\t NAME\t DRIVER\t")

	for _, n := range networks {
		fmt.Fprintf(w, "%s\t %s\t %s\t\n",
			truncateID(n.ID),
			truncate(n.Name, 20),
			truncate(n.Driver, 20),
		)
	}
	w.Flush()
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

func truncateID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}
