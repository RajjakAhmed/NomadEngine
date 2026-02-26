package workflow

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/RajjakAhmed/NomadEngine/internal/store"
)

func CreateStepExecution(stepID string) (string, time.Time, error) {
	executionID := uuid.New().String()
	startedAt := time.Now()

	_, err := store.DB.Exec(context.Background(),
		`INSERT INTO step_executions 
		 (id, step_id, status, started_at)
		 VALUES ($1, $2, $3, $4)`,
		executionID,
		stepID,
		StepRunning,
		startedAt,
	)

	return executionID, startedAt, err
}

func CompleteStepExecution(executionID string, success bool, errMsg string, startedAt time.Time) error {
	completedAt := time.Now()
	duration := completedAt.Sub(startedAt).Milliseconds()

	status := StepCompleted
	if !success {
		status = StepFailed
	}

	_, err := store.DB.Exec(context.Background(),
		`UPDATE step_executions 
		 SET status=$1,
		     error=$2,
		     completed_at=$3,
		     duration_ms=$4
		 WHERE id=$5`,
		status,
		errMsg,
		completedAt,
		duration,
		executionID,
	)

	return err
}