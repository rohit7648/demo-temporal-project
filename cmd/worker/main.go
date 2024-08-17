package main

import (
	"log"
	"os"
	"os/signal"

	"demo-temporal-project/client/temporal"
	"demo-temporal-project/configs"
	"demo-temporal-project/constant"
	"demo-temporal-project/internal/usecase/bank"
	"demo-temporal-project/pkg/logger"

	"go.temporal.io/sdk/worker"
)

// @@@SNIPSTART money-transfer-project-template-go-worker
func main() {

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)

	config, err := configs.NewConfig()
	if err != nil {
		log.Fatalf("Unable to load config")
	}

	logger := logger.NewLogger(config.Log)
	c, cleanup, err := temporal.NewTemporalClient(logger, config.Temporal, config.Env)

	if err != nil {
		log.Fatalf("Unable to dial Temporal client, err - %s", err)
	}
	defer cleanup()

	w := worker.New(*c, constant.MoneyTransferTaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflow(bank.MoneyTransfer)
	w.RegisterActivity(bank.Withdraw)
	w.RegisterActivity(bank.Deposit)
	w.RegisterActivity(bank.Refund)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
	<-stopCh
}

// @@@SNIPEND
