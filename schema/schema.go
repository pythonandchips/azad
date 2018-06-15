package schema

import "github.com/pythonandchips/azad/conn"

// Schema schema
type Schema struct {
	Tasks  []Task
	Fields []Field
}

// Task executable
type Task struct {
	Name   string
	Fields []Field
	Run    func(map[string]string, conn.Conn) error
}

// Field field
type Field struct {
	Name     string
	Type     string
	Required bool
}

func BuildSchema() Schema {
	return Schema{}
}
