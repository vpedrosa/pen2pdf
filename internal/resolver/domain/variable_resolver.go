package domain

import (
	"fmt"
	"strings"

	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// VariableResolver walks the document tree and replaces $variable references
// with their concrete values from the document's variables map.
type VariableResolver struct{}

func NewVariableResolver() *VariableResolver {
	return &VariableResolver{}
}

func (r *VariableResolver) Resolve(doc *shared.Document) error {
	if doc.Variables == nil {
		return nil
	}
	for _, child := range doc.Children {
		if err := resolveNode(child, doc.Variables); err != nil {
			return err
		}
	}
	return nil
}

func resolveNode(node shared.Node, vars map[string]shared.Variable) error {
	switch n := node.(type) {
	case *shared.Frame:
		return resolveFrame(n, vars)
	case *shared.Text:
		return resolveText(n, vars)
	default:
		return fmt.Errorf("unsupported node type: %T", node)
	}
}

func resolveFrame(frame *shared.Frame, vars map[string]shared.Variable) error {
	if frame.Fill != nil && frame.Fill.Type == shared.FillSolid {
		resolved, err := resolveColorString(frame.Fill.Color, vars)
		if err != nil {
			return fmt.Errorf("frame %q fill: %w", frame.ID, err)
		}
		frame.Fill.Color = resolved
	}

	for _, child := range frame.Children {
		if err := resolveNode(child, vars); err != nil {
			return err
		}
	}
	return nil
}

func resolveText(text *shared.Text, vars map[string]shared.Variable) error {
	resolved, err := resolveColorString(text.Fill, vars)
	if err != nil {
		return fmt.Errorf("text %q fill: %w", text.ID, err)
	}
	text.Fill = resolved
	return nil
}

func resolveColorString(value string, vars map[string]shared.Variable) (string, error) {
	if !strings.HasPrefix(value, "$") {
		return value, nil
	}

	name := value[1:]
	v, ok := vars[name]
	if !ok {
		return "", fmt.Errorf("undefined variable: %q", name)
	}

	str, ok := v.Value.(string)
	if !ok {
		return "", fmt.Errorf("variable %q is not a string (type: %s)", name, v.Type)
	}
	return str, nil
}
