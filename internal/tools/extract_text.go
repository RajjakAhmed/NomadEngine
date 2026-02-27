package tools

import "fmt"

func init() {
	Register(Tool{
		Name: "extract_text",
		Run:  extractText,
	})
}

func extractText(input map[string]interface{}) (map[string]interface{}, error) {

	memory, ok := input["memory"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("memory not found")
	}

	fetch, ok := memory["fetch_url"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("fetch_url output not found")
	}

	text, ok := fetch["text"].(string)
	if !ok {
		return nil, fmt.Errorf("text not found in fetch_url")
	}

	// For now just return the same text
	return map[string]interface{}{
		"text": text,
	}, nil
}