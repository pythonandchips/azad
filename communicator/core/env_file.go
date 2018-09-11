package core

import (
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pythonandchips/azad/plugin"
)

func envFile(context plugin.InputContext) (map[string]string, error) {
	variables := map[string]string{}
	filePath := context.Get("path")
	absolutePath := filepath.Join(context.PlaybookPath(), filePath)
	vars, err := godotenv.Read(absolutePath)
	if err != nil {
		return variables, err
	}
	for key, value := range vars {
		variables[strings.ToLower(key)] = value
	}
	return variables, nil
}
