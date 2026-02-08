package domain

type Document struct {
	Version   string
	Children  []Node
	Variables map[string]Variable
}
