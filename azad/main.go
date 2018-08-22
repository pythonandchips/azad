package main

import (
	"log"
	"os"

	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/logger"
	"github.com/urfave/cli"
)

func main() {
	logger.Initialize()

	app := app()
	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func app() *cli.App {
	app := cli.NewApp()
	app.Name = "Azad: Server Configuration Management"
	app.Flags = flags()
	app.Action = actionHander(runPlaybook)
	return app
}

func flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "user, u",
			Value: "root",
			Usage: "user for ssh connection",
		},
		cli.StringFlag{
			Name:  "key, k",
			Value: conn.DefaultSSHKeyPath(),
			Usage: "ssh key used to connect to server",
		},
		cli.BoolFlag{
			Name:  "simulate, s",
			Usage: "Simulate run and output script to stdout, used for development",
		},
	}
}
