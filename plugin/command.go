package plugin

// Command represents a command to be ran on a remote host
type Command struct {
	// Interpreter used to run command e.g. sh, bash, ruby
	Interpreter string
	// Command lines used to make up a command to be ran
	Command []string
	// Additional environment variable to be used in the command
	Env map[string]string
}
