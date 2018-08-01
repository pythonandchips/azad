package core

import (
	"bytes"
	"fmt"
	"path/filepath"

	tmpl "text/template"

	"github.com/pythonandchips/azad/plugin"
)

func templateConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "source", Type: "String", Required: true},
			{Name: "dest", Type: "String", Required: true},
		},
		Run: template,
	}
}

func template(context plugin.Context) (map[string]string, error) {
	templatePath := filepath.Join(context.RolePath(), "templates", context.Get("source"))
	fileTemplate, _ := tmpl.New(context.Get("dest")).ParseFiles(templatePath)
	buf := bytes.NewBuffer([]byte{})
	fileTemplate.Execute(buf, context.Variables())
	command := plugin.Command{
		Interpreter: "sh",
		Command: []string{
			"filecont=<<- HEREDOC",
			buf.String(),
			"HEREDOC",
			fmt.Sprintf("echo $filecount > %s", context.Get("dest")),
		},
	}
	err := context.Run(command)
	return map[string]string{}, err
}
