package service

import (
	"context"
	pb "demo-temporal-project/api-protos/demo_temporal_project/v1"
	"demo-temporal-project/configs"
	"demo-temporal-project/constant"
	"demo-temporal-project/internal/usecase/bank"

	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransactionService struct {
	pb.UnimplementedTransactionServer
	TemporalClient *client.Client
}

func NewTransactionServer(tc *client.Client) *TransactionService {

	return &TransactionService{
		TemporalClient: tc,
	}
}

func (ts *TransactionService) TransferMoney(ctx context.Context, req *pb.PaymentDetails) (*emptypb.Empty, error) {

	tc := *ts.TemporalClient
	_, err := tc.ExecuteWorkflow(ctx, getWorkflowOptions(req, constant.MoneyTransferTaskQueueName), bank.MoneyTransfer, &pb.PaymentDetails{
		SourceAccount: req.SourceAccount,
		TargetAccount: req.TargetAccount,
		Amount:        req.Amount,
		RefId:         req.RefId,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func getWorkflowOptions(req *pb.PaymentDetails, taskQueue string) client.StartWorkflowOptions {
	wflowCfg := configs.Cfg.Temporal.Workflow
	return client.StartWorkflowOptions{
		ID:                  req.RefId,
		TaskQueue:           taskQueue,
		WorkflowTaskTimeout: wflowCfg.TaskTimeout,
	}
}
