package azad

import (
	"fmt"
	"strings"
)

// Inventory info required to generate servers from a provider e.g. aws
type Inventory struct {
	Name       string
	Attributes map[string]string
}

// PluginName returns the plugin name for the inventory
func (inventory Inventory) PluginName() (string, error) {
	parts := strings.Split(inventory.Name, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid inventory name %s", inventory.Name)
	}
	return parts[0], nil
}

// ServiceName returns the service name for the inventory
func (inventory Inventory) ServiceName() (string, error) {
	parts := strings.Split(inventory.Name, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid inventory name %s", inventory.Name)
	}
	return parts[1], nil
}
