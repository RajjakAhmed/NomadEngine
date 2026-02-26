package tools

import "fmt"

type ToolFunc func(input map[string]interface{}) (map[string]interface{}, error)

var registry = make(map[string]ToolFunc)

func Register(name string, fn ToolFunc) {
	registry[name] = fn
}

func Execute(name string, input map[string]interface{}) (map[string]interface{}, error) {
	tool, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	return tool(input)
}