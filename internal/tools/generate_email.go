package tools

func generateEmail(input map[string]interface{}) (map[string]interface{}, error) {

	context, _ := input["context"].(string)

	email := "Hello,\n\n" +
		"Based on our research:\n\n" +
		context +
		"\n\nBest,\nNomad Engine"

	return map[string]interface{}{
		"email": email,
	}, nil
}

func init() {
	Register(Tool{
		Name:        "generate_email",
		Description: "Generate professional outreach email",
		Run:         generateEmail,
	})
}