package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"

	shared "github.com/kannanmohan1994/common-structs"

	temporal_client "go.temporal.io/sdk/client"
)

const (
	TaskQueueAssignDriverWorkflow = "TASK_QUEUE_ASSIGN_DRIVER_WORKFLOW"
)

func main() {
	c, err := client.NewClient(client.Options{
		HostPort:  "localhost:7233",
		Namespace: "default",
	})
	if err != nil {
		log.Fatalln("unable to connect temporal client", err)
	}
	defer c.Close()

	execute(c)
}

func execute(temporalClient client.Client) error {
	workflowOptions := temporal_client.StartWorkflowOptions{
		TaskQueue: TaskQueueAssignDriverWorkflow,
		ID:        "assign-driver-" + uuid.New().String(),
	}
	_, err := temporalClient.ExecuteWorkflow(
		context.Background(),
		workflowOptions,
		"AssignDriverWorkflow",
		&shared.AssignDriverWorkflowRequest{
			UserID:   "123",
			Location: "12345",
		},
	)
	return err
}
