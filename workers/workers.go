package workers

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"log"
	"temporal-demo/models"
	"temporal-demo/utils"
	_workflows "temporal-demo/workflows"
	"time"
)

// ExecuteStep Activity functions
func SampleActivity(ctx context.Context, step models.Activities) (string, error) {
	return fmt.Sprintf("Executed step: %s with task: %s", step.Name, step.Parameters), nil
}

// DynamicWorkflow Workflow function
func DynamicWorkflow(ctx workflow.Context, parameters map[string]interface{}) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	workflows, _ := _workflows.LoadWorkflows()

	wrkflow, _ := utils.FindWorkflowByID(workflows, "321")

	var result string
	for _, step := range wrkflow.Activities {
		err := workflow.ExecuteActivity(ctx, step.Name, step.Parameters).Get(ctx, &result)
		if err != nil {
			return "", err
		}
	}

	return "Workflow completed", nil
}

// StartTemporalWorker Start Temporal worker
func StartTemporalWorker() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic("Failed to create Temporal client")
	}
	defer c.Close()

	w := worker.New(c, "dynamic-workflow-task-queue", worker.Options{})

	workflows, err := _workflows.LoadWorkflows()
	if err != nil {
		log.Fatal(err)
	}
	wrkflow, err := utils.FindWorkflowByID(workflows, "321")
	if err != nil {
		log.Fatal(err)
	}
	//
	//function, err := activities.GetFuncByName(wrkflow.ID)
	w.RegisterWorkflow(wrkflow.ID)
	for _, step := range wrkflow.Activities {
		w.RegisterActivity(step.Name)
	}

	err = w.Run(worker.InterruptCh())
	if err != nil {
		panic("Failed to start Temporal worker")
	}
}
