package bank

import (
	"context"
	pb "demo-temporal-project/api-protos/demo_temporal_project/v1"
	app "demo-temporal-project/client/banking"
	"fmt"
	"log"
)

func Withdraw(ctx context.Context, data *pb.PaymentDetails) (string, error) {
	log.Printf("Withdrawing $%f from account %s.\n\n",
		data.Amount,
		data.SourceAccount,
	)

	referenceID := fmt.Sprintf("%s-withdrawal", data.RefId)
	bank := app.BankingService{Hostname: "bank-api.example.com"}
	confirmation, err := bank.Withdraw(data.SourceAccount, int(data.Amount), referenceID)
	return confirmation, err
}

func Deposit(ctx context.Context, data *pb.PaymentDetails) (string, error) {
	log.Printf("Depositing $%f into account %s.\n\n",
		data.Amount,
		data.TargetAccount,
	)

	referenceID := fmt.Sprintf("%s-deposit", data.RefId)
	bank := app.BankingService{Hostname: "bank-api.example.com"}
	// Uncomment the next line and comment the one after that to simulate an unknown failure
	// confirmation, err := bank.DepositThatFails(data.TargetAccount, data.Amount, referenceID)
	confirmation, err := bank.Deposit(data.TargetAccount, int(data.Amount), referenceID)
	return confirmation, err
}

func Refund(ctx context.Context, data *pb.PaymentDetails) (string, error) {
	log.Printf("Refunding $%v back into account %v.\n\n",
		data.Amount,
		data.SourceAccount,
	)

	referenceID := fmt.Sprintf("%s-refund", data.RefId)
	bank := app.BankingService{Hostname: "bank-api.example.com"}
	confirmation, err := bank.Deposit(data.SourceAccount, int(data.Amount), referenceID)
	return confirmation, err
}
