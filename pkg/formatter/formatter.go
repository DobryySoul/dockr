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

func Confirm(question string, resources *domain.UnusedResources) bool {
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

func PrintReport(res *domain.UnusedResources, dryRun bool) {
	if dryRun {
		color.New(color.FgYellow).Println("\n=== DRY RUN MODE ===")
		color.New(color.FgHiBlue).Println("Будут удалены следующие ресурсы:")
	} else {
		color.New(color.FgYellow).Println("\n=== ОЧИСТКА DOCKER ===")
	}

	printSection := func(title string, count int, fn func()) {
		if count > 0 {
			color.New(color.FgGreen).Printf("\n%s (%d):\n", title, count)
			fn()
		}
	}

	printSection("Образы", len(res.Images), func() {
		printImagesTable(res.Images)
	})

	printSection("Контейнеры", len(res.Containers), func() {
		printContainersTable(res.Containers)
	})

	// printSection("Тома", len(res.Volumes), func() {
	// 	printVolumesTable(res.Volumes)
	// })

	// printSection("Сети", len(res.Networks), func() {
	// 	printNetworksTable(res.Networks)
	// })

	if res.TotalCount() > 0 {
		color.New(color.FgHiWhite).Printf("\nИтого: %d ресурсов, ", res.TotalCount())
		color.New(color.FgHiGreen).Printf("освободится %.2f MB\n", res.TotalSize()/1024/1024)
	} else {
		color.New(color.FgHiGreen).Println("\nНечего удалять - система чиста!")
	}
}

func printImagesTable(images []image.Summary) {
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

func printContainersTable(containers []container.Summary) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\t NAME\t СОСТОЯНИЕ\t IMAGE\t")

	for _, c := range containers {
		state := c.State
		if c.State == "exited" {
			state = color.HiRedString(state)
		}

		fmt.Fprintf(w, "%s\t %s\t %s\t %s\t\n",
			truncateID(c.ID),
			truncate(c.Names[0], 20),
			state,
			truncate(c.Image, 20),
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
