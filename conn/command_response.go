package conn

import "bytes"

// CommandResponse wraps output from stdout and stderr
type CommandResponse struct {
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

// Stdout returns anything written to stdout while running a command
func (commandResposne CommandResponse) Stdout() string {
	return commandResposne.stdout.String()
}

// Stderr returns anything written to stderr while running a command
func (commandResposne CommandResponse) Stderr() string {
	return commandResposne.stderr.String()
}
