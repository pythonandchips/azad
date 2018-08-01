package conn

import (
	"fmt"
	"strings"
)

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

func (command Command) generateFile() string {
	interpreter := command.Interpreter
	if interpreter == "" {
		interpreter = "sh"
	}
	fileLines := []string{
		fmt.Sprintf("#!/usr/bin/env %s", interpreter),
		"",
	}
	for key, value := range command.Env {
		fileLines = append(fileLines, fmt.Sprintf("%s='%s'", key, value))
	}
	fileLines = append(fileLines, "")
	for _, line := range command.Command {
		fileLines = append(fileLines, line)
	}
	return strings.Join(fileLines, "\n")
}
