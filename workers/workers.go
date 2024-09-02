package workers

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	dsl "temporal-demo/activities"
	_workflows "temporal-demo/workflows"
)

// StartTemporalWorker Start Temporal worker
func StartTemporalWorker() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic("Failed to create Temporal client")
	}
	defer c.Close()

	// The client and worker are heavyweight objects that should be created once per process.
	wl := worker.New(c, "dsl-workflow-task-queue", worker.Options{})

	wl.RegisterWorkflow(_workflows.SimpleDSLWorkflow)
	wl.RegisterActivity(&dsl.SampleActivities{})

	err = wl.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

}
