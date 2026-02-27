package engine

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/RajjakAhmed/NomadEngine/internal/planner"
	"github.com/RajjakAhmed/NomadEngine/internal/workflow"
)

type GoalRequest struct {
	Goal string `json:"goal"`
}

func RegisterRoutes(router *gin.Engine) {
	router.POST("/workflows/:id/execute", executeWorkflowHandler)
	router.POST("/autonomous", autonomousHandler)
}

func autonomousHandler(c *gin.Context) {
	var req GoalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// create workflow
	wf, err := workflow.CreateWorkflow(req.Goal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// planner generates steps
	steps := planner.Plan(req.Goal)

	for i, step := range steps {
		_, err := workflow.AddStep(wf.ID, step.Action, step.Input, i+1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// run executor
	go ExecuteWorkflow(wf.ID)

	c.JSON(http.StatusOK, gin.H{
		"workflow_id": wf.ID,
		"message":     "autonomous workflow started",
	})
}

func executeWorkflowHandler(c *gin.Context) {
	id := c.Param("id")

	err := workflow.UpdateWorkflowStatus(id, workflow.StatusRunning)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func() {
		ExecuteWorkflow(id)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "workflow execution started"})
}