package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"temporal-demo/models"
	"temporal-demo/utils"
	"time"

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

func ExecuteDSLWorkFlow(c *gin.Context) {

	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dslWorkflows, err := LoadWorkflows()
	dslWorkflow, err := utils.FindWorkflowByID(dslWorkflows, input["name"])
	dslWorkflow.Variables = input

	// The client is a heavyweight object that should be created once per process.
	co, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        "dsl_" + uuid.New().String(),
		TaskQueue: "dsl-workflow-task-queue",
	}

	we, err := co.ExecuteWorkflow(context.Background(), workflowOptions, SimpleDSLWorkflow, *dslWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	c.JSON(http.StatusOK, "Started workflow"+we.GetID())

}

func startTemporalWorkflow(workflow models.Workflow) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic("Failed to create Temporal client")
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		TaskQueue: "dynamic-workflow-task-queue",
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, workflow.Variables)
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

// SimpleDSLWorkflow workflow definition
func SimpleDSLWorkflow(ctx workflow.Context, dslWorkflow models.Workflow) ([]byte, error) {
	bindings := make(map[string]string)
	//workflowcheck:ignore Only iterates for building another map
	for k, v := range dslWorkflow.Variables {
		bindings[k] = v
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	err := dslWorkflow.Root.Execute(ctx, bindings)
	if err != nil {
		logger.Error("DSL Workflow failed.", "Error", err)
		return nil, err
	}

	logger.Info("DSL Workflow completed.")
	return nil, err
}
