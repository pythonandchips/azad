package core

import (
	"bytes"
	"path"
	"path/filepath"

	"text/template"

	"github.com/pythonandchips/azad/plugin"
	"github.com/pythonandchips/azad/plugin/helpers"
)

func templateConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "source", Type: "String", Required: true},
			{Name: "dest", Type: "String", Required: true},
			{Name: "locals", Type: "Map", Required: false},
		},
		Run: templateCommand,
	}
}

func templateCommand(context plugin.Context) (map[string]string, error) {
	templatePath := filepath.Join(context.RolePath(), "templates", context.Get("source"))
	templateName := path.Base(context.Get("source"))
	funcMap := template.FuncMap{
		"get": func(key string) string {
			return context.Get(key)
		},
	}
	fileTemplate, err := template.New(templateName).Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		return map[string]string{}, err
	}
	data := bytes.NewBuffer([]byte{})
	err = fileTemplate.Execute(data, context.GetMap("locals"))
	if err != nil {
		return map[string]string{}, err
	}
	commands := []string{}
	commands = append(commands, helpers.Checksum(data.Bytes(), context.Get("dest"))...)
	commands = append(commands, helpers.WriteEncodedFile(data.Bytes(), context.Get("dest"))...)
	command := plugin.Command{
		Command: commands,
	}
	err = context.Run(command)
	return map[string]string{}, err
}
