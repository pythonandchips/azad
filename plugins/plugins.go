package plugins

import (
	"os"
	"path/filepath"
	"plugin"

	"github.com/pythonandchips/azad/core"
	"github.com/pythonandchips/azad/schema"
)

var pluginLoader PluginLoader

func Configure() {
	pluginLoader = GoPluginLoader{}
}

// Loader loader
func Loader() PluginLoader {
	return pluginLoader
}

// PluginLoader plugin loader
type PluginLoader interface {
	Get(string) (Plugin, error)
	TaskExists(string, string) error
	GetTask(string, string) (schema.Task, error)
}

// GoPluginLoader go plugin loader
type GoPluginLoader struct {
	loadPath string
}

func (p GoPluginLoader) load(pluginName string) (Plugin, error) {
	plugin, err := p.loadPlugin(pluginName)
	if err != nil {
		return goPlugin{}, err
	}
	schemaFn, err := plugin.Lookup("Schema")
	if err != nil {
		return goPlugin{}, err
	}
	schema := schemaFn.(func() schema.Schema)()
	pluginList[pluginName] = goPlugin{
		plugin: plugin,
		schema: schema,
		name:   pluginName,
	}
	return pluginList[pluginName], nil
}

func (p GoPluginLoader) loadPlugin(pluginName string) (IPlugin, error) {
	if pluginName == "core" {
		return core.New(), nil
	}
	loadPath := ".plugins"
	if p.loadPath == "" {
		loadPath = p.loadPath
	}
	dir, _ := os.Getwd()
	pluginPath := filepath.Join(dir, loadPath, pluginName+".so")
	return plugin.Open(pluginPath)
}

// Get get
func (p GoPluginLoader) Get(pluginName string) (Plugin, error) {
	plugin, found := pluginList[pluginName]
	if found {
		return plugin, nil
	}
	return p.load(pluginName)
}

// TaskExists check task
func (p GoPluginLoader) TaskExists(pluginName, taskName string) error {
	plugin, err := p.Get(pluginName)
	if err != nil {
		return err
	}
	return plugin.TaskExists(taskName)
}

// GetTask form plugin
func (p GoPluginLoader) GetTask(pluginName, taskName string) (schema.Task, error) {
	plugin, err := p.Get(pluginName)
	if err != nil {
		return schema.Task{}, err
	}
	return plugin.GetTask(taskName)
}
