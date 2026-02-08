package domain

const (
	NodeTypeFrame = "frame"
	NodeTypeText  = "text"
)

type Node interface {
	GetID() string
	GetName() string
	GetType() string
}

type Frame struct {
	ID             string
	Name           string
	X              float64
	Y              float64
	Width          Dimension
	Height         Dimension
	Fill           *Fill
	CornerRadius   float64
	Clip           bool
	Layout         string
	Gap            float64
	Padding        Padding
	JustifyContent string
	AlignItems     string
	Children       []Node
}

func (f *Frame) GetID() string   { return f.ID }
func (f *Frame) GetName() string { return f.Name }
func (f *Frame) GetType() string { return NodeTypeFrame }

type Text struct {
	ID            string
	Name          string
	Content       string
	Fill          string
	FontFamily    string
	FontSize      float64
	FontWeight    string
	FontStyle     string
	LetterSpacing float64
	LineHeight    float64
	TextAlign     string
	Width         Dimension
	TextGrowth    string
}

func (t *Text) GetID() string   { return t.ID }
func (t *Text) GetName() string { return t.Name }
func (t *Text) GetType() string { return NodeTypeText }

type Dimension struct {
	Value         float64
	FillContainer bool
}

func FixedDimension(v float64) Dimension {
	return Dimension{Value: v}
}

func FillContainerDimension() Dimension {
	return Dimension{FillContainer: true}
}

type Padding struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

func UniformPadding(v float64) Padding {
	return Padding{Top: v, Right: v, Bottom: v, Left: v}
}
