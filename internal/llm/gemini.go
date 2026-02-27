package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GeminiClient struct {
	APIKey string
}

func NewGemini(apiKey string) *GeminiClient {
	return &GeminiClient{
		APIKey: apiKey,
	}
}

func (g *GeminiClient) Generate(prompt string) (string, error) {

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s",
		g.APIKey,
	)

	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Gemini raw response:")
	fmt.Println(string(respBody))

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("gemini API error: %s", string(respBody))
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", err
	}

	candidatesRaw, ok := result["candidates"]
	if !ok {
		return "", fmt.Errorf("gemini response missing candidates")
	}

	candidates, ok := candidatesRaw.([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("invalid candidates format")
	}

	firstCandidate, ok := candidates[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid candidate structure")
	}

	contentRaw, ok := firstCandidate["content"]
	if !ok {
		return "", fmt.Errorf("missing content field")
	}

	contentMap, ok := contentRaw.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid content structure")
	}

	partsRaw, ok := contentMap["parts"]
	if !ok {
		return "", fmt.Errorf("missing parts field")
	}

	parts, ok := partsRaw.([]interface{})
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("invalid parts structure")
	}

	partMap, ok := parts[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid part structure")
	}

	textRaw, ok := partMap["text"]
	if !ok {
		return "", fmt.Errorf("missing text field")
	}

	text, ok := textRaw.(string)
	if !ok {
		return "", fmt.Errorf("invalid text format")
	}

	return text, nil
}