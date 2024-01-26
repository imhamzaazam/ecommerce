package workflow

import (
	pb "ecommerce/codegen/transactions/v1/proto"
	"ecommerce/transactions/store"
	"go.temporal.io/sdk/workflow"
	"time"
)

func CreateTransactionWorkflow(ctx workflow.Context, request *pb.CreateTransactionRequest) (*pb.Transaction, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	transaction := store.ProtoToStruct(request)
	response := pb.CreateTransactionResponse{}
	err := workflow.ExecuteActivity(ctx, transaction).Get(ctx, &response)
	var transactionResponse pb.Transaction
	transactionResponse.Id = response.Id
	transactionResponse.CreatedAt = response.CreatedAt
	return &transactionResponse, err
}
