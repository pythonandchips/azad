package main

import (
	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/runner"
)

func runPlaybook(c context) error {
	if c.Bool("simulate") {
		logger.Info("Simulating run, no change will be made to server")
		conn.Simulate()
	}
	playbookFilePath := "./playbook.az"
	if c.NArg() > 0 {
		playbookFilePath = c.Args().Get(0)
	}
	logger.Info("Applying %s with user %s", playbookFilePath, c.String("user"))
	runner.RunPlaybook(playbookFilePath, runner.Config{
		KeyFilePath: c.String("key"),
		User:        c.String("user"),
	})
	return nil
}
