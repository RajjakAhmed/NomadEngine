package tools

func webSearch(input map[string]interface{}) (map[string]interface{}, error) {

	query, _ := input["query"].(string)

	results := []map[string]string{
		{
			"name": "Stripe",
			"url":  "https://stripe.com",
		},
		{
			"name": "Plaid",
			"url":  "https://plaid.com",
		},
		{
			"name": "Ramp",
			"url":  "https://ramp.com",
		},
	}

	return map[string]interface{}{
		"query":   query,
		"results": results,
	}, nil
}