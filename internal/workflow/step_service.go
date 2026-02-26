package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/RajjakAhmed/NomadEngine/internal/store"
)

// AddStep adds a new step to a workflow with the given action and input
func AddStep(workflowID, action string, input map[string]interface{}, order int) (*Step, error) {
	// Check workflow status
	var wfStatus string
	err := store.DB.QueryRow(context.Background(),
		`SELECT status FROM workflows WHERE id=$1`,
		workflowID,
	).Scan(&wfStatus)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("workflow not found")
		}
		return nil, fmt.Errorf("failed to check workflow status: %w", err)
	}

	if wfStatus != StatusPending {
		return nil, fmt.Errorf("cannot add steps unless workflow is pending")
	}

	// Check if step order already exists for this workflow
	var count int
	err = store.DB.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM steps WHERE workflow_id=$1 AND step_order=$2`,
		workflowID, order,
	).Scan(&count)

	if err != nil {
		return nil, fmt.Errorf("failed to check step order: %w", err)
	}

	if count > 0 {
		return nil, fmt.Errorf("step order %d already exists for this workflow", order)
	}

	// Create the step
	id := uuid.New().String()
	now := time.Now()

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	query := `
		INSERT INTO steps (id, workflow_id, step_order, action, input, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = store.DB.Exec(context.Background(), query,
		id, workflowID, order, action, inputJSON, StatusPending, now, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create step: %w", err)
	}

	return &Step{
		ID:         id,
		WorkflowID: workflowID,
		StepOrder:  order,
		Action:     action,
		Input:      input,
		Status:     StatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// GetStepsByWorkflow retrieves all steps for a given workflow, ordered by step_order
func GetStepsByWorkflow(workflowID string) ([]Step, error) {
	query := `
		SELECT id, workflow_id, step_order, action, input, status, retry_count, max_retries, created_at, updated_at
		FROM steps
		WHERE workflow_id = $1
		ORDER BY step_order ASC
	`

	rows, err := store.DB.Query(context.Background(), query, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to query steps: %w", err)
	}
	defer rows.Close()

	var steps []Step

	for rows.Next() {
		var step Step
		var inputBytes []byte

		err := rows.Scan(
			&step.ID,
			&step.WorkflowID,
			&step.StepOrder,
			&step.Action,
			&inputBytes,
			&step.Status,
			&step.RetryCount,
			&step.MaxRetries,
			&step.CreatedAt,
			&step.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan step: %w", err)
		}

		err = json.Unmarshal(inputBytes, &step.Input)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal input: %w", err)
		}

		steps = append(steps, step)
	}

	// Check for errors after iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return steps, nil
}

// UpdateStepStatus updates the status of a step by its ID
func UpdateStepStatus(stepID string, newStatus string) error {
	result, err := store.DB.Exec(context.Background(),
		`UPDATE steps SET status=$1, updated_at=$2 WHERE id=$3`,
		newStatus,
		time.Now(),
		stepID,
	)

	if err != nil {
		return fmt.Errorf("failed to update step status: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("step not found")
	}

	return nil
}
// IncrementRetry increments the retry count for a step and updates the updated_at timestamp
func IncrementRetry(stepID string) error {
	_, err := store.DB.Exec(context.Background(),
		`UPDATE steps 
		 SET retry_count = retry_count + 1,
		     updated_at = $1
		 WHERE id = $2`,
		time.Now(),
		stepID,
	)
	return err
}