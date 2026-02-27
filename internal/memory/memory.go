package memory

import (
	"context"
	"encoding/json"
	"github.com/RajjakAhmed/NomadEngine/internal/store"
)

func Save(workflowID, key string, value map[string]interface{}) error {

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO workflow_memory (workflow_id, key, value)
	VALUES ($1, $2, $3)
	`

	_, err = store.DB.Exec(context.Background(), query, workflowID, key, data)
	return err
}

func LoadAll(workflowID string) (map[string]interface{}, error) {

	query := `
	SELECT key, value
	FROM workflow_memory
	WHERE workflow_id = $1
	`

	rows, err := store.DB.Query(context.Background(), query, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]interface{}{}

	for rows.Next() {

		var key string
		var raw []byte

		if err := rows.Scan(&key, &raw); err != nil {
			return nil, err
		}

		var value interface{}
		json.Unmarshal(raw, &value)

		result[key] = value
	}

	return result, nil
}