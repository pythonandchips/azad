package main

import (
	"log"
	"os"

	"github.com/pythonandchips/azad/azad"
	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/plugins"
	"github.com/pythonandchips/azad/runner"
	"github.com/urfave/cli"
)

func main() {
	plugins.Configure()
	logger.Initialize()

	app := cli.NewApp()
	app.Name = "Azad: Server Configuration Management"
	app.Flags = flags()
	app.Action = runPlaybook

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "user, u",
			Value: "root",
			Usage: "user for ssh connection",
		},
		cli.BoolFlag{
			Name:  "simulate, s",
			Usage: "Simulate run and output script to stdout, used for development",
		},
	}
}

func runPlaybook(c *cli.Context) error {
	if c.Bool("simulate") {
		logger.Info("Simulating run, no change will be made to server")
		conn.Simulate()
	}
	playbookFilePath := "./playbook.az"
	if c.NArg() > 0 {
		playbookFilePath = c.Args().Get(0)
	}
	logger.Info("Applying %s", playbookFilePath)
	runner.RunPlaybook(playbookFilePath, azad.Config{})
	return nil
}
