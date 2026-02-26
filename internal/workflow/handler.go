package workflow

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	
)

type CreateWorkflowRequest struct {
	Goal string `json:"goal" binding:"required"`
}

// RegisterRoutes sets up the API routes for workflow management

func RegisterRoutes(router *gin.Engine) {
	router.POST("/workflows", createWorkflowHandler)
	router.GET("/workflows/:id", getWorkflowHandler)

	router.PATCH("/workflows/:id/start", startWorkflowHandler)
	router.PATCH("/workflows/:id/complete", completeWorkflowHandler)
	router.PATCH("/workflows/:id/fail", failWorkflowHandler)
	router.POST("/workflows/:id/steps", addStepHandler)
	router.GET("/workflows/:id/steps", getStepsHandler)
	
}

func createWorkflowHandler(c *gin.Context) {
	var req CreateWorkflowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflow, err := CreateWorkflow(req.Goal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to create workflow: %v", err),
		})
		return
	}

	c.JSON(http.StatusCreated, workflow)
}

func getWorkflowHandler(c *gin.Context) {
	id := c.Param("id")

	workflow, err := GetWorkflowByID(id)

	// Handle "not found" case (workflow is nil, err is nil)
	if workflow == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "workflow not found",
		})
		return
	}

	// Handle actual database errors
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get workflow: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, workflow)
}
func startWorkflowHandler(c *gin.Context) {
	id := c.Param("id")

	err := UpdateWorkflowStatus(id, StatusRunning)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workflow started"})
}

func completeWorkflowHandler(c *gin.Context) {
	id := c.Param("id")

	err := UpdateWorkflowStatus(id, StatusCompleted)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workflow completed"})
}

func failWorkflowHandler(c *gin.Context) {
	id := c.Param("id")

	err := UpdateWorkflowStatus(id, StatusFailed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workflow failed"})
}

// create step request struct
type CreateStepRequest struct {
	Action string                 `json:"action" binding:"required"`
	Input  map[string]interface{} `json:"input" binding:"required"`
	Order  int                    `json:"order" binding:"required"`
}

func addStepHandler(c *gin.Context) {
	workflowID := c.Param("id")

	var req CreateStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	step, err := AddStep(workflowID, req.Action, req.Input, req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, step)
}

func getStepsHandler(c *gin.Context) {
	workflowID := c.Param("id")

	steps, err := GetStepsByWorkflow(workflowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, steps)
}
