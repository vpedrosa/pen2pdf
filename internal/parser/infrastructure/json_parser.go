package infrastructure

import (
	"encoding/json"
	"fmt"
	"io"

	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// JSONParser parses .pen files (JSON format) into the domain model.
type JSONParser struct{}

func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

func (p *JSONParser) Parse(r io.Reader) (*shared.Document, error) {
	var raw rawDocument
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	children, err := parseNodes(raw.Children)
	if err != nil {
		return nil, err
	}

	variables, err := parseVariables(raw.Variables)
	if err != nil {
		return nil, err
	}

	return &shared.Document{
		Version:   raw.Version,
		Children:  children,
		Variables: variables,
	}, nil
}

// rawDocument is the top-level JSON structure.
type rawDocument struct {
	Version   string                     `json:"version"`
	Children  []json.RawMessage          `json:"children"`
	Variables map[string]json.RawMessage `json:"variables"`
}

// rawNode is a partially-decoded node used to determine type.
type rawNode struct {
	Type string `json:"type"`
}

// rawFrame holds all frame fields with polymorphic types as RawMessage.
type rawFrame struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	X              float64           `json:"x"`
	Y              float64           `json:"y"`
	Width          json.RawMessage   `json:"width"`
	Height         json.RawMessage   `json:"height"`
	Fill           json.RawMessage   `json:"fill"`
	CornerRadius   float64           `json:"cornerRadius"`
	Clip           bool              `json:"clip"`
	Layout         string            `json:"layout"`
	Gap            float64           `json:"gap"`
	Padding        json.RawMessage   `json:"padding"`
	JustifyContent string            `json:"justifyContent"`
	AlignItems     string            `json:"alignItems"`
	Children       []json.RawMessage `json:"children"`
}

// rawText holds all text fields with polymorphic types as RawMessage.
type rawText struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Content       string          `json:"content"`
	Fill          string          `json:"fill"`
	FontFamily    string          `json:"fontFamily"`
	FontSize      float64         `json:"fontSize"`
	FontWeight    string          `json:"fontWeight"`
	FontStyle     string          `json:"fontStyle"`
	LetterSpacing float64         `json:"letterSpacing"`
	LineHeight    float64         `json:"lineHeight"`
	TextAlign     string          `json:"textAlign"`
	Width         json.RawMessage `json:"width"`
	TextGrowth    string          `json:"textGrowth"`
}

func parseNodes(rawNodes []json.RawMessage) ([]shared.Node, error) {
	nodes := make([]shared.Node, 0, len(rawNodes))
	for i, raw := range rawNodes {
		node, err := parseNode(raw)
		if err != nil {
			return nil, fmt.Errorf("child[%d]: %w", i, err)
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func parseNode(data json.RawMessage) (shared.Node, error) {
	var probe rawNode
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("cannot determine node type: %w", err)
	}

	switch probe.Type {
	case "frame":
		return parseFrame(data)
	case "text":
		return parseText(data)
	default:
		return nil, fmt.Errorf("unknown node type: %q", probe.Type)
	}
}

func parseFrame(data json.RawMessage) (*shared.Frame, error) {
	var raw rawFrame
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid frame: %w", err)
	}

	width, err := parseDimension(raw.Width)
	if err != nil {
		return nil, fmt.Errorf("frame %q width: %w", raw.ID, err)
	}

	height, err := parseDimension(raw.Height)
	if err != nil {
		return nil, fmt.Errorf("frame %q height: %w", raw.ID, err)
	}

	fill, err := parseFill(raw.Fill)
	if err != nil {
		return nil, fmt.Errorf("frame %q fill: %w", raw.ID, err)
	}

	padding, err := parsePadding(raw.Padding)
	if err != nil {
		return nil, fmt.Errorf("frame %q padding: %w", raw.ID, err)
	}

	children, err := parseNodes(raw.Children)
	if err != nil {
		return nil, fmt.Errorf("frame %q: %w", raw.ID, err)
	}

	return &shared.Frame{
		ID:             raw.ID,
		Name:           raw.Name,
		X:              raw.X,
		Y:              raw.Y,
		Width:          width,
		Height:         height,
		Fill:           fill,
		CornerRadius:   raw.CornerRadius,
		Clip:           raw.Clip,
		Layout:         raw.Layout,
		Gap:            raw.Gap,
		Padding:        padding,
		JustifyContent: raw.JustifyContent,
		AlignItems:     raw.AlignItems,
		Children:       children,
	}, nil
}

func parseText(data json.RawMessage) (*shared.Text, error) {
	var raw rawText
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid text: %w", err)
	}

	width, err := parseDimension(raw.Width)
	if err != nil {
		return nil, fmt.Errorf("text %q width: %w", raw.ID, err)
	}

	return &shared.Text{
		ID:            raw.ID,
		Name:          raw.Name,
		Content:       raw.Content,
		Fill:          raw.Fill,
		FontFamily:    raw.FontFamily,
		FontSize:      raw.FontSize,
		FontWeight:    raw.FontWeight,
		FontStyle:     raw.FontStyle,
		LetterSpacing: raw.LetterSpacing,
		LineHeight:    raw.LineHeight,
		TextAlign:     raw.TextAlign,
		Width:         width,
		TextGrowth:    raw.TextGrowth,
	}, nil
}

// parseDimension handles: number (800), string ("fill_container"), or absent (null/empty).
func parseDimension(data json.RawMessage) (shared.Dimension, error) {
	if len(data) == 0 || string(data) == "null" {
		return shared.Dimension{}, nil
	}

	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		return shared.FixedDimension(num), nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "fill_container" {
			return shared.FillContainerDimension(), nil
		}
		return shared.Dimension{}, fmt.Errorf("unknown dimension value: %q", str)
	}

	return shared.Dimension{}, fmt.Errorf("invalid dimension: %s", string(data))
}

// parseFill handles: string ("#RRGGBB"), object ({type, url, ...}), or absent.
func parseFill(data json.RawMessage) (*shared.Fill, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		return shared.SolidFill(str), nil
	}

	var obj struct {
		Type    string  `json:"type"`
		URL     string  `json:"url"`
		Mode    string  `json:"mode"`
		Opacity float64 `json:"opacity"`
		Enabled bool    `json:"enabled"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("invalid fill: %w", err)
	}

	switch obj.Type {
	case "image":
		return shared.ImageFill(obj.URL, obj.Mode, obj.Opacity, obj.Enabled), nil
	default:
		return nil, fmt.Errorf("unknown fill type: %q", obj.Type)
	}
}

// parsePadding handles: number (40), 2-element array [v, h], 4-element array [t, r, b, l], or absent.
func parsePadding(data json.RawMessage) (shared.Padding, error) {
	if len(data) == 0 || string(data) == "null" {
		return shared.Padding{}, nil
	}

	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		return shared.UniformPadding(num), nil
	}

	var arr []float64
	if err := json.Unmarshal(data, &arr); err != nil {
		return shared.Padding{}, fmt.Errorf("invalid padding: %s", string(data))
	}

	switch len(arr) {
	case 2:
		return shared.Padding{Top: arr[0], Right: arr[1], Bottom: arr[0], Left: arr[1]}, nil
	case 4:
		return shared.Padding{Top: arr[0], Right: arr[1], Bottom: arr[2], Left: arr[3]}, nil
	default:
		return shared.Padding{}, fmt.Errorf("padding array must have 2 or 4 elements, got %d", len(arr))
	}
}

func parseVariables(rawVars map[string]json.RawMessage) (map[string]shared.Variable, error) {
	if rawVars == nil {
		return nil, nil
	}

	variables := make(map[string]shared.Variable, len(rawVars))
	for name, data := range rawVars {
		v, err := parseVariable(data)
		if err != nil {
			return nil, fmt.Errorf("variable %q: %w", name, err)
		}
		variables[name] = v
	}
	return variables, nil
}

func parseVariable(data json.RawMessage) (shared.Variable, error) {
	var raw struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return shared.Variable{}, fmt.Errorf("invalid variable: %w", err)
	}

	varType := shared.VariableType(raw.Type)
	switch varType {
	case shared.VariableColor, shared.VariableString:
		var s string
		if err := json.Unmarshal(raw.Value, &s); err != nil {
			return shared.Variable{}, fmt.Errorf("expected string value: %w", err)
		}
		return shared.Variable{Type: varType, Value: s}, nil
	case shared.VariableNumber:
		var n float64
		if err := json.Unmarshal(raw.Value, &n); err != nil {
			return shared.Variable{}, fmt.Errorf("expected number value: %w", err)
		}
		return shared.Variable{Type: varType, Value: n}, nil
	default:
		return shared.Variable{}, fmt.Errorf("unknown variable type: %q", raw.Type)
	}
}
