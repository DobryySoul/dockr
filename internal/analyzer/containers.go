package analyzer

import "github.com/docker/docker/api/types/container"

// IsContainerUnused проверяет, является ли контейнер неиспользуемым (остановленным).
// Контейнеры со статусом "exited", "created", "dead" можно безопасно удалять.
func IsContainerUnused(c *container.Summary) bool {
	// Допустимые статусы для удаления
	return c.State == "exited" || c.State == "created" || c.State == "dead"
}
