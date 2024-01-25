package store

import (
	pb "ecommerce/codegen/transactions/v1/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ProtoToStruct(in *pb.CreateTransactionRequest) Transaction {
	return Transaction{
		Id:         uuid.MustParse(in.Id),
		CustomerId: uuid.MustParse(in.CustomerId),
		ProductId:  uuid.MustParse(in.ProductId),
		Quantity:   in.Quantity,
		Amount:     in.Amount,
		CreatedAt: timestamppb.Timestamp{
			Seconds: int64(time.Now().Second()),
			Nanos:   int32(time.Now().Nanosecond()),
		},
	}
}

func StructToProto(transaction *Transaction) *pb.Transaction {
	return &pb.Transaction{
		Id:         transaction.Id.String(),
		CustomerId: transaction.CustomerId.String(),
		ProductId:  transaction.ProductId.String(),
		Quantity:   transaction.Quantity,
		Amount:     transaction.Amount,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: transaction.CreatedAt.Seconds,
			Nanos:   transaction.CreatedAt.Nanos,
		},
	}
}
