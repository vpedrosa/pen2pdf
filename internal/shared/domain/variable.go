package domain

type VariableType string

const (
	VariableColor  VariableType = "color"
	VariableString VariableType = "string"
	VariableNumber VariableType = "number"
)

type Variable struct {
	Type  VariableType
	Value any
}
