package engine

import (
	"fmt"

	"github.com/RajjakAhmed/NomadEngine/internal/workflow"
	"github.com/RajjakAhmed/NomadEngine/internal/tools"
)

func ExecuteWorkflow(workflowID string) error {
	// Get all steps for the workflow
	steps, err := workflow.GetStepsByWorkflow(workflowID)
	if err != nil {
		return fmt.Errorf("failed to get workflow steps: %w", err)
	}

	// Execute each step with retry logic
	for _, step := range steps {
		// Mark step as running
		err = workflow.UpdateStepStatus(step.ID, workflow.StepRunning)
		if err != nil {
			return fmt.Errorf("failed to mark step as running: %w", err)
		}

		success := false

		// Retry loop with execution tracking
		for step.RetryCount < step.MaxRetries {
			// Create step execution record
			executionID, startedAt, err := workflow.CreateStepExecution(step.ID)
			if err != nil {
				return fmt.Errorf("failed to create step execution: %w", err)
			}

			fmt.Printf("Executing step: %s (Attempt: %d/%d) [Execution ID: %s]\n",
				step.Action,
				step.RetryCount+1,
				step.MaxRetries,
				executionID,
			)

			// Execute the tool
			output, err := tools.Execute(step.Action, step.Input)

			if err != nil {
				// Increment retry count
				incrementErr := workflow.IncrementRetry(step.ID)
				if incrementErr != nil {
					return fmt.Errorf("failed to increment retry count: %w", incrementErr)
				}
				step.RetryCount++

				// Record failed execution
				completeErr := workflow.CompleteStepExecution(
					executionID,
					false,
					err.Error(),
					startedAt,
				)
				if completeErr != nil {
					return fmt.Errorf("failed to record execution failure: %w", completeErr)
				}

				fmt.Printf("Step %s failed: %v\n", step.Action, err)
				continue
			}

			// Record successful execution
			completeErr := workflow.CompleteStepExecution(
				executionID,
				true,
				"",
				startedAt,
			)
			if completeErr != nil {
				return fmt.Errorf("failed to record execution success: %w", completeErr)
			}

			fmt.Printf("Step %s output: %v\n", step.Action, output)

			// Step succeeded
			success = true
			break
		}

		// Check if step failed after all retries
		if !success {
			fmt.Printf("Step %s failed after %d retries\n", step.Action, step.MaxRetries)

			err = workflow.UpdateStepStatus(step.ID, workflow.StepFailed)
			if err != nil {
				return fmt.Errorf("failed to update step status: %w", err)
			}

			err = workflow.UpdateWorkflowStatus(workflowID, workflow.StatusFailed)
			if err != nil {
				return fmt.Errorf("failed to update workflow status: %w", err)
			}

			return fmt.Errorf("step %s failed after max retries", step.Action)
		}

		// Mark step as completed
		err = workflow.UpdateStepStatus(step.ID, workflow.StepCompleted)
		if err != nil {
			return fmt.Errorf("failed to mark step as completed: %w", err)
		}

		fmt.Printf("Step %s completed successfully\n", step.Action)
	}

	// Mark workflow as completed
	err = workflow.UpdateWorkflowStatus(workflowID, workflow.StatusCompleted)
	if err != nil {
		return fmt.Errorf("failed to mark workflow as completed: %w", err)
	}

	fmt.Printf("Workflow %s completed successfully\n", workflowID)
	return nil
}