package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"temporal-demo/models"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

const workflowsFile = "workflows.json"

func CreateWorkflow(c *gin.Context) {
	var workflow models.Workflow
	if err := c.ShouldBindJSON(&workflow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflows, err := LoadWorkflows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	workflows = append(workflows, workflow)
	if err := saveWorkflows(workflows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Start the workflow in Temporal
	go startTemporalWorkflow(workflow)

	c.JSON(http.StatusOK, workflow)
}

func LoadWorkflows() ([]models.Workflow, error) {
	var workflows []models.Workflow

	data, err := ioutil.ReadFile(workflowsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return workflows, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &workflows); err != nil {
		return nil, err
	}

	return workflows, nil
}

func saveWorkflows(workflows []models.Workflow) error {
	data, err := json.MarshalIndent(workflows, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(workflowsFile, data, 0644)
}

func startTemporalWorkflow(workflow models.Workflow) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic("Failed to create Temporal client")
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        workflow.ID,
		TaskQueue: "dynamic-workflow-task-queue",
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, workflow.ID, workflow.Parameters)
	if err != nil {
		panic("Failed to start workflow")
	}

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		panic("Failed to get workflow result")
	}

	fmt.Printf("Workflow result: %s\n", result)
}
