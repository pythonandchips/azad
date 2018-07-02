package plugin

// Schema schema
type Schema struct {
	Tasks     map[string]Task
	Inventory map[string]Inventory
}

// Task executable
type Task struct {
	Fields []Field
	Run    func(Context) error
}

// Field field
type Field struct {
	Name     string
	Type     string
	Required bool
}

// Inventory represents how to host information from an external system
type Inventory struct {
	Fields []Field
	Run    func(Context) ([]Resource, error)
}

// Resource represents a potential host with the host name for the server and
// Groups it will be a member of
type Resource struct {
	ConnectOn string
	Groups    []string
}
