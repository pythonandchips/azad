package main

import "github.com/urfave/cli"

type context interface {
	Bool(string) bool
	Args() cli.Args
	NArg() int
	String(string) string
}

type cliContext struct {
	context *cli.Context
}

func (context cliContext) Bool(key string) bool {
	return context.context.Bool(key)
}

func (context cliContext) Args() cli.Args {
	return context.context.Args()
}

func (context cliContext) NArg() int {
	return context.context.NArg()
}

func (context cliContext) String(key string) string {
	return context.context.String(key)
}
