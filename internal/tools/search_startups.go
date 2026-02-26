package tools
import (
	"fmt"
)

func init() {
	Register("search_startups", SearchStartups)
}

func SearchStartups(input map[string]interface{}) (map[string]interface{}, error) {
	industry := input["industry"]
	limit := input["limit"]

	result := map[string]interface{}{
		"results": fmt.Sprintf("Top %v startups in %v industry", limit, industry),
	}

	return result, nil
}