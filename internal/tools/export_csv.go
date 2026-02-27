package tools

import (
	"encoding/csv"
	"os"
)

func exportCSV(input map[string]interface{}) (map[string]interface{}, error) {

	rows, ok := input["rows"].([][]string)
	if !ok {
		return nil, nil
	}

	file, err := os.Create("export.csv")
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)
	writer.WriteAll(rows)

	writer.Flush()

	return map[string]interface{}{
		"file": "export.csv",
	}, nil
}

func init() {
	Register(Tool{
		Name:        "export_csv",
		Description: "Export structured data to CSV file",
		Run:         exportCSV,
	})
}
