package main

import "github.com/urfave/cli"

func actionHander(f func(context) error) func(*cli.Context) error {
	return func(c *cli.Context) error {
		return f(cliContext{c})
	}
}
