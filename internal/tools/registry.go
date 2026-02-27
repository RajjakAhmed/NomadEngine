package tools

import "fmt"

type Tool struct {
	Name        string
	Description string
	Run         func(input map[string]interface{}) (map[string]interface{}, error)
}

var registry = make(map[string]Tool)

func Register(tool Tool) {
	registry[tool.Name] = tool
}

func Execute(name string, input map[string]interface{}) (map[string]interface{}, error) {

	tool, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	return tool.Run(input)
}

func ListTools() []Tool {

	var tools []Tool

	for _, t := range registry {
		tools = append(tools, t)
	}

	return tools
}