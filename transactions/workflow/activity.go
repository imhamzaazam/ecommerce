package workflow

import (
	"context"
	pb "ecommerce/codegen/transactions/v1/proto"
	"ecommerce/transactions/store"
	v1 "ecommerce/transactions/v1"
)

type Activity struct {
	transactionServiceClient v1.TransactionService
}

func (a *Activity) CreateTransaction(ctx context.Context, request *store.Transaction) (*pb.CreateTransactionResponse, error) {
	return a.transactionServiceClient.CreateTransaction(ctx, &pb.CreateTransactionRequest{
		Id:         request.Id.String(),
		CustomerId: request.CustomerId.String(),
		ProductId:  request.ProductId.String(),
		Quantity:   request.Quantity,
		Amount:     request.Amount,
	})
}
