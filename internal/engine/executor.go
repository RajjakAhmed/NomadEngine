package engine

import (
	"fmt"

	"github.com/RajjakAhmed/NomadEngine/internal/memory"
	"github.com/RajjakAhmed/NomadEngine/internal/tools"
	"github.com/RajjakAhmed/NomadEngine/internal/workflow"
)

func ExecuteWorkflow(workflowID string) error {

	// Get all steps
	steps, err := workflow.GetStepsByWorkflow(workflowID)
	if err != nil {
		return fmt.Errorf("failed to get workflow steps: %w", err)
	}

	for _, step := range steps {

		err = workflow.UpdateStepStatus(step.ID, workflow.StepRunning)
		if err != nil {
			return fmt.Errorf("failed to mark step running: %w", err)
		}

		success := false

		for step.RetryCount < step.MaxRetries {

			executionID, startedAt, err := workflow.CreateStepExecution(step.ID)
			if err != nil {
				return fmt.Errorf("failed to create execution record: %w", err)
			}

			fmt.Printf(
				"Executing step: %s (Attempt %d/%d) [Execution ID: %s]\n",
				step.Action,
				step.RetryCount+1,
				step.MaxRetries,
				executionID,
			)

			//----------------------------------
			// Load workflow memory
			//----------------------------------

			mem, err := memory.LoadAll(workflowID)
			if err != nil {
				fmt.Printf("Warning: memory load failed: %v\n", err)
				mem = map[string]interface{}{}
			}

			//----------------------------------
			// Prepare input safely
			//----------------------------------

			enrichedInput := map[string]interface{}{}

			// copy step input
			for k, v := range step.Input {
				enrichedInput[k] = v
			}

			// inject memory
			enrichedInput["memory"] = mem

			fmt.Printf("Injected memory: %v\n", mem)

			//----------------------------------
			// Execute tool
			//----------------------------------

			output, err := tools.Execute(step.Action, enrichedInput)

			if err != nil {

				workflow.IncrementRetry(step.ID)
				step.RetryCount++

				workflow.CompleteStepExecution(
					executionID,
					false,
					err.Error(),
					startedAt,
				)

				fmt.Printf("Step %s failed: %v\n", step.Action, err)
				continue
			}

			//----------------------------------
			// Save output to memory
			//----------------------------------

			err = memory.Save(workflowID, step.Action, output)
			if err != nil {
				fmt.Printf("Warning: memory save failed: %v\n", err)
			}

			workflow.CompleteStepExecution(
				executionID,
				true,
				"",
				startedAt,
			)

			fmt.Printf("Step %s output: %v\n", step.Action, output)

			success = true
			break
		}

		if !success {

			fmt.Printf("Step %s failed after %d retries\n", step.Action, step.MaxRetries)

			workflow.UpdateStepStatus(step.ID, workflow.StepFailed)
			workflow.UpdateWorkflowStatus(workflowID, workflow.StatusFailed)

			return fmt.Errorf("step %s failed after max retries", step.Action)
		}

		err = workflow.UpdateStepStatus(step.ID, workflow.StepCompleted)
		if err != nil {
			return fmt.Errorf("failed to mark step completed: %w", err)
		}

		fmt.Printf("Step %s completed successfully\n", step.Action)
	}

	err = workflow.UpdateWorkflowStatus(workflowID, workflow.StatusCompleted)
	if err != nil {
		return fmt.Errorf("failed to mark workflow completed: %w", err)
	}

	fmt.Printf("Workflow %s completed successfully\n", workflowID)

	return nil
}