package conn

import (
	"fmt"
	"strings"
)

type Command struct {
	Interpreter string
	Command     []string
	Env         map[string]string
	User        string
}

func (command Command) generateFile() string {
	fileLines := []string{
		fmt.Sprintf("#!/usr/bin/env %s", command.Interpreter),
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
