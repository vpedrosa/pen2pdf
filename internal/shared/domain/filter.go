package domain

import (
	"fmt"
	"strings"
)

// FilterPagesByName filters a slice of Nodes, keeping only those whose
// GetName() is in the comma-separated names list.
func FilterPagesByName(children []Node, names string) ([]Node, error) {
	nameList := strings.Split(names, ",")
	nameSet := make(map[string]bool, len(nameList))
	for _, n := range nameList {
		nameSet[strings.TrimSpace(n)] = true
	}

	var filtered []Node
	for _, child := range children {
		if nameSet[child.GetName()] {
			filtered = append(filtered, child)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no pages match filter %q", names)
	}
	return filtered, nil
}
