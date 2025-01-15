package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"userservice/helpers"

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
	verifyUserWorker := worker.New(temporalClient, shared.TaskQueueVerifyUserActivity, worker.Options{})
	verifyUserWorker.RegisterActivity(helpers.VerifyUserActivity)

	fetchUserDetailsWorker := worker.New(temporalClient, shared.TaskQueueFetchUserDetailsActivity, worker.Options{})
	fetchUserDetailsWorker.RegisterActivity(helpers.FetchUserDetailsActivity)

	return []worker.Worker{verifyUserWorker, fetchUserDetailsWorker}
}
