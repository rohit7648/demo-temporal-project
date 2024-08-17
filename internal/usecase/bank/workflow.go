package bank

import (
	"fmt"
	"time"

	pb "demo-temporal-project/api-protos/demo_temporal_project/v1"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func MoneyTransfer(ctx workflow.Context, input *pb.PaymentDetails) (string, error) {
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        500,
		NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retrypolicy,
	}

	// Apply the options to the context
	ctx = workflow.WithActivityOptions(ctx, options)

	var withdrawOutput string
	// Execute the Withdraw activity
	withdrawErr := workflow.ExecuteActivity(ctx, Withdraw, input).Get(ctx, &withdrawOutput)
	if withdrawErr != nil {
		return "", withdrawErr
	}

	var depositOutput string
	// Execute the Deposit activity
	depositErr := workflow.ExecuteActivity(ctx, Deposit, input).Get(ctx, &depositOutput)
	if depositErr != nil {
		var result string
		// Execute the Refund activity if Deposit fails
		refundErr := workflow.ExecuteActivity(ctx, Refund, input).Get(ctx, &result)
		if refundErr != nil {
			return "", fmt.Errorf("Deposit: failed to deposit money into %v: %v. Money could not be returned to %v: %w",
				input.TargetAccount, depositErr, input.SourceAccount, refundErr)
		}
		return "", fmt.Errorf("Deposit: failed to deposit money into %v: Money returned to %v: %w",
			input.TargetAccount, input.SourceAccount, depositErr)
	}

	result := fmt.Sprintf("Transfer complete (transaction IDs: %s, %s)", withdrawOutput, depositOutput)
	return result, nil
}
