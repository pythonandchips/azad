package aptget

import (
	"fmt"
	"regexp"

	"github.com/pythonandchips/azad/plugin"
)

func addRepoConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "repo", Type: "String", Required: true},
			{Name: "key", Type: "String", Required: true},
			{Name: "update", Type: "Bool"},
		},
		Run: addRepoCommand,
	}
}

func addRepoCommand(context plugin.Context) (map[string]string, error) {
	// regexp removes the url from the debian mapage name
	// expected format is
	// "deb https://some.url.to.packages/sub-but arch main"
	// and "some.url.to.packages" will be extracted
	packageRegex, err := regexp.Compile(`(?:deb https:\/\/)(\S*)(?:\/.*)`)
	if err != nil {
		return map[string]string{}, err
	}
	matches := packageRegex.FindAllStringSubmatch(context.Get("repo"), -1)
	if len(matches) != 1 {
		return map[string]string{}, fmt.Errorf("unable to find package name")
	}
	basename := matches[0][1]
	commands := []string{
		fmt.Sprintf(`wget -qO - %s | sudo apt-key add -`, context.Get("key")),
		fmt.Sprintf(`echo "%s" > /etc/apt/sources.list.d/%s.list`, context.Get("repo"), basename),
	}
	if context.IsTrue("update") {
		commands = append(commands, "apt-get update -qy")
	}
	command := plugin.Command{
		Command: commands,
	}
	err = context.Run(command)
	return map[string]string{}, err
}
