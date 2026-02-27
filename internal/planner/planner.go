package planner

import (
	"log"
	"os"
	"strings"

	"github.com/RajjakAhmed/NomadEngine/internal/llm"
)

type PlannedStep struct {
	Action string                 `json:"action"`
	Input  map[string]interface{} `json:"input"`
}

func Plan(goal string) []PlannedStep {

	apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey != "" {

		client := llm.NewGemini(apiKey)

		steps, err := LLMPlan(goal, client)
		if err == nil && len(steps) > 0 {
			log.Println("Planner: using Gemini LLM")
			return steps
		}

		log.Println("Gemini planner error:", err)
		log.Println("Planner: Gemini failed, falling back to rule planner")
	}

	return rulePlan(goal)
}

func rulePlan(goal string) []PlannedStep {

	goal = strings.ToLower(goal)

	var steps []PlannedStep

	if strings.Contains(goal, "startup") {
		steps = append(steps, PlannedStep{
			Action: "search_startups",
			Input: map[string]interface{}{
				"query": "fintech startups",
			},
		})
	}

	if strings.Contains(goal, "email") {
		steps = append(steps, PlannedStep{
			Action: "draft_email",
			Input:  map[string]interface{}{},
		})
	}

	return steps
}