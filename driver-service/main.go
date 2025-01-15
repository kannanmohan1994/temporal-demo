package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"userservice/helpers"

	temporal_activity "go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	shared "github.com/kannanmohan1994/common-structs"
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
	matchDriverWorker := worker.New(temporalClient, shared.TaskQueueMatchDriverActivity, worker.Options{})
	matchDriverWorker.RegisterActivityWithOptions(helpers.MatchDriverActivity, temporal_activity.RegisterOptions{
		Name: shared.MatchDriverActivity,
	})
	waitDriverWorker := worker.New(temporalClient, shared.TaskQueueWaitDriverActivity, worker.Options{})
	waitDriverWorker.RegisterActivityWithOptions(helpers.WaitDriverActivity, temporal_activity.RegisterOptions{
		Name: shared.WaitDriverActivity,
	})

	return []worker.Worker{matchDriverWorker, waitDriverWorker}
}
