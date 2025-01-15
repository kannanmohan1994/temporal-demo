package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"workflows/helpers"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
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

	var blocker = make(chan os.Signal, 1)
	signal.Notify(blocker, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	for _, temporalWorker := range initWorkers(c) {
		err = temporalWorker.Start()
		if err != nil {
			log.Fatalln("unable to start worker", err)
		}
	}
	<-blocker
}

func initWorkers(temporalClient client.Client) []worker.Worker {
	driverWorkflow := helpers.Driverworkflow{Client: temporalClient}
	rideWorkflow := worker.New(temporalClient, TaskQueueAssignDriverWorkflow, worker.Options{})
	rideWorkflow.RegisterWorkflow(driverWorkflow.AssignDriverWorkflow)

	return []worker.Worker{rideWorkflow}
}
