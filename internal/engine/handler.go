package engine

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/RajjakAhmed/NomadEngine/internal/workflow"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/workflows/:id/execute", executeWorkflowHandler)
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