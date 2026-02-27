package tools

import (
	"fmt"
	"strings"
)

func init() {
	Register(Tool{
		Name: "summarize_text",
		Run:  summarizeText,
	})
}

func summarizeText(input map[string]interface{}) (map[string]interface{}, error) {

	memory, ok := input["memory"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("memory not found")
	}

	extract, ok := memory["extract_text"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("extract_text output not found")
	}

	text, ok := extract["text"].(string)
	if !ok {
		return nil, fmt.Errorf("text not found")
	}

	// Simple summary logic (first 400 chars)
	text = strings.TrimSpace(text)

	if len(text) > 400 {
		text = text[:400] + "..."
	}

	return map[string]interface{}{
		"summary": text,
	}, nil
}