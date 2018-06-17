package schema

// Schema schema
type Schema struct {
	Tasks  []Task
	Fields []Field
}

// Task executable
type Task struct {
	Name   string
	Fields []Field
	Run    func(Context) error
}

// Field field
type Field struct {
	Name     string
	Type     string
	Required bool
}

// Command represents a command to be ran on a remote host
type Command struct {
	// Interpreter used to run command e.g. sh, bash, ruby
	Interpreter string
	// Command lines used to make up a command to be ran
	Command []string
	// Additional environment variable to be used in the command
	Env map[string]string
	// User to run the command with e.g. root, admin
	User string
}
