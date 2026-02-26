package workflow

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/RajjakAhmed/NomadEngine/internal/store"
)

func CreateWorkflow(goal string) (*Workflow, error) {
	if goal == "" {
		return nil, fmt.Errorf("goal cannot be empty")
	}

	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO workflows (id, goal, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := store.DB.Exec(context.Background(), query,
		id, goal, "pending", now, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	return &Workflow{
		ID:        id,
		Goal:      goal,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func GetWorkflowByID(id string) (*Workflow, error) {
	query := `
		SELECT id, goal, status, created_at, updated_at, started_at, completed_at
		FROM workflows
		WHERE id = $1
	`

	row := store.DB.QueryRow(context.Background(), query, id)

	var wf Workflow
	var startedAt, completedAt sql.NullTime

	err := row.Scan(
		&wf.ID,
		&wf.Goal,
		&wf.Status,
		&wf.CreatedAt,
		&wf.UpdatedAt,
		&startedAt,
		&completedAt,
	)

	if err != nil {
		// Distinguish between "not found" and actual database errors
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil, nil for "not found"
		}
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	// Convert sql.NullTime to *time.Time
	if startedAt.Valid {
		wf.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		wf.CompletedAt = &completedAt.Time
	}

	return &wf, nil
}
func UpdateWorkflowStatus(id string, newStatus string) error {
	ctx := context.Background()

	// Get current status
	var currentStatus string
	err := store.DB.QueryRow(ctx,
		`SELECT status FROM workflows WHERE id = $1`,
		id,
	).Scan(&currentStatus)

	if err != nil {
		return err
	}

	// Validate transitions
	if !isValidTransition(currentStatus, newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
	}

	// Update status
	_, err = store.DB.Exec(ctx,
		`UPDATE workflows SET status=$1, updated_at=$2 WHERE id=$3`,
		newStatus,
		time.Now(),
		id,
	)

	return err
}
func isValidTransition(current, next string) bool {
	switch current {
	case StatusPending:
		return next == StatusRunning
	case StatusRunning:
		return next == StatusCompleted || next == StatusFailed
	default:
		return false
	}
}