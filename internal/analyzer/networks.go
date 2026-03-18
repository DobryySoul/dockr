package analyzer

import "github.com/docker/docker/api/types/network"

// IsNetworkUnused проверяет, является ли сеть неиспользуемой.
// В Docker сеть считается неиспользуемой, если к ней не подключен ни один контейнер.
// Базовые сети Docker (bridge, host, none) обычно не следует удалять.
func IsNetworkUnused(net *network.Summary) bool {
	// Игнорируем стандартные сети Docker
	if net.Name == "bridge" || net.Name == "host" || net.Name == "none" {
		return false
	}
	// Если к сети не подключено ни одного контейнера (Containers map пустая)
	return len(net.Containers) == 0
}
