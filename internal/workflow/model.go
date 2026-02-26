package workflow

import "time"
const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)
type Workflow struct {
	ID          string    `json:"id"`
	Goal        string    `json:"goal"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
const (
	StepPending   = "pending"
	StepRunning   = "running"
	StepCompleted = "completed"
	StepFailed    = "failed"
)
type Step struct {
	ID         string                 `json:"id"`
	WorkflowID string                 `json:"workflow_id"`
	StepOrder  int                    `json:"step_order"`
	Action     string                 `json:"action"`
	Input      map[string]interface{} `json:"input"`
	Status     string                 `json:"status"`
	RetryCount int                    `json:"retry_count"`
	MaxRetries int 					  `json:"max_retries"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}