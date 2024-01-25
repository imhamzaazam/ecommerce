package v1

import (
	"context"
	pb "ecommerce/codegen/transactions/v1/proto"
	"ecommerce/transactions/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

type TransactionService struct {
	pb.UnimplementedTransactionsServer
	transactionChannel chan pb.Transaction
	transactionStore   *store.TransactionStore
}

func NewTransactionService(conn *pgx.Conn) *TransactionService {
	return &TransactionService{transactionChannel: make(chan pb.Transaction), transactionStore: store.NewTransactionStore(conn)}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, in *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	log.Printf("Received transaction request: %v", in.Id)

	transaction := store.ProtoToStruct(in)
	err := s.transactionStore.InsertTransaction(ctx, &transaction)
	if err != nil {
		return nil, err
	}

	select {

	case s.transactionChannel <- pb.Transaction{
		Id:         transaction.Id.String(),
		CustomerId: transaction.CustomerId.String(),
		ProductId:  transaction.ProductId.String(),
		Quantity:   transaction.Quantity,
		Amount:     transaction.Amount,
		CreatedAt: &timestamppb.Timestamp{Seconds: transaction.CreatedAt.Seconds,
			Nanos: transaction.CreatedAt.Nanos,
		},
	}:
	default:
		log.Println("testing")
	}

	return &pb.CreateTransactionResponse{
		Id:        in.Id,
		CreatedAt: timestamppb.Now(),
	}, nil

}

func (s *TransactionService) GetTransactionById(ctx context.Context, in *pb.GetTransactionRequest) (
	*pb.GetTxnResponse,
	error,
) {
	log.Printf("Recieved Get transaction by id request: %v", in.GetId())
	//fetch from db
	transaction, err := s.transactionStore.GetTransactionById(ctx, uuid.MustParse(in.Id))
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, status.Error(codes.NotFound, "transaction not found!")
	}
	return &pb.GetTxnResponse{
		Request: store.StructToProto(transaction),
	}, nil
}

func (s *TransactionService) ListTransactions(in *pb.ListTransactionRequest, srv pb.Transactions_ListTransactionsServer) error {
	log.Printf("Recieved: ")
	for {
		transaction := <-s.transactionChannel
		if err := srv.Send(&pb.GetTxnResponse{
			Request: &pb.Transaction{
				Id:        transaction.Id,
				ProductId: transaction.ProductId,
			},
		}); err != nil {
			log.Println("error generating response")
			return err
		}
	}
}
