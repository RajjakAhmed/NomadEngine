package planner

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RajjakAhmed/NomadEngine/internal/llm"
	"github.com/RajjakAhmed/NomadEngine/internal/tools"
)

func LLMPlan(goal string, client *llm.GeminiClient) ([]PlannedStep, error) {

	// ----------------------------
	// Get tool metadata
	// ----------------------------

	toolList := tools.ListTools()

	toolDescriptions := ""

	for _, t := range toolList {
		toolDescriptions += fmt.Sprintf("%s: %s\n", t.Name, t.Description)
	}

	// ----------------------------
	// Build prompt
	// ----------------------------

	prompt := fmt.Sprintf(`
You are a workflow planning AI.

Convert the user goal into workflow steps.

Goal:
%s

Available tools:
%s

Rules:
- Use only the tools listed above
- Each step must contain "action" and "input"
- Return ONLY valid JSON
- Do NOT include explanations

Example format:

[
  {
    "action": "web_search",
    "input": {"query": "fintech startups"}
  },
  {
    "action": "summarize_text",
    "input": {"text": "search results"}
  }
]
`, goal, toolDescriptions)

	// ----------------------------
	// Call Gemini
	// ----------------------------

	resp, err := client.Generate(prompt)
	if err != nil {
		return nil, err
	}

	// ----------------------------
	// Clean markdown (LLMs often wrap JSON)
	// ----------------------------

	resp = strings.TrimSpace(resp)
	resp = strings.TrimPrefix(resp, "```json")
	resp = strings.TrimPrefix(resp, "```")
	resp = strings.TrimSuffix(resp, "```")

	// ----------------------------
	// Parse JSON
	// ----------------------------

	var steps []PlannedStep

	err = json.Unmarshal([]byte(resp), &steps)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM JSON: %v\nresponse: %s", err, resp)
	}

	if len(steps) == 0 {
		return nil, fmt.Errorf("LLM returned empty workflow")
	}

	return steps, nil
}