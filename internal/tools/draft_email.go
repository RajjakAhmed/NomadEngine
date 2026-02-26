package tools
func init() {
	Register("draft_email", DraftEmail)
}

func DraftEmail(input map[string]interface{}) (map[string]interface{}, error) {
	tone := input["tone"]

	return map[string]interface{}{
		"email": "This is a " + tone.(string) + " outreach email draft.",
	}, nil
}