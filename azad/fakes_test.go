package main

import "github.com/urfave/cli"

type FakeContext struct {
	bools   map[string]bool
	strings map[string]string
	args    []string
}

func (fakeContext FakeContext) Bool(key string) bool {
	return fakeContext.bools[key]
}

func (fakeContext FakeContext) Args() cli.Args {
	return fakeContext.args
}

func (fakeContext FakeContext) NArg() int {
	return len(fakeContext.args)
}

func (fakeContext FakeContext) String(key string) string {
	return fakeContext.strings[key]
}
