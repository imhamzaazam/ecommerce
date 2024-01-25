package store

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type TransactionStore struct {
	session *pgx.Conn
}
type Transaction struct {
	Id         uuid.UUID
	CustomerId uuid.UUID
	ProductId  uuid.UUID
	Quantity   int32
	Amount     float64
	CreatedAt  timestamppb.Timestamp
}

func NewTransactionStore(session *pgx.Conn) *TransactionStore {
	return &TransactionStore{session: session}
}

func (store *TransactionStore) InsertTransaction(ctx context.Context, trx *Transaction) error {
	// Insert transactions in transaction table.
	log.Println("Inserting into transactions table")
	t := timestamppb.Timestamp{Seconds: trx.CreatedAt.Seconds,
		Nanos: trx.CreatedAt.Nanos,
	}
	//var exists bool
	//row := store.session.QueryRow(ctx,
	//	"SELECT EXISTS(SELECT 1 FROM transactions where id = $1)", trx.Id)
	//if err := row.Scan(&exists); err != nil {
	//	return err
	//}
	//if !exists {
	if _, err2 := store.session.Exec(ctx,
		"INSERT INTO transactions (id, customer_id, product_id, quantity, amount, created) VALUES ($1, $2, $3, $4, $5, $6)",
		trx.Id, trx.CustomerId, trx.ProductId, trx.Quantity, trx.Amount, t.AsTime()); err2 != nil {
		return err2
	}
	//}
	return nil
}

func (store *TransactionStore) GetTransactionById(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	// Insert transactions in transaction table.
	log.Println("Get from transactions table by ID")
	var created time.Time
	var transaction = &Transaction{}
	err := store.session.QueryRow(ctx,
		"SELECT id, product_id, customer_id, quantity, amount, created FROM transactions WHERE id = $1", id).Scan(
		&transaction.Id,
		&transaction.ProductId,
		&transaction.CustomerId,
		&transaction.Quantity,
		&transaction.Amount,
		&created,
	)
	if err != nil {
		return nil, err
	}
	transaction.CreatedAt = timestamppb.Timestamp{Seconds: int64(created.Second()),
		Nanos: int32(created.Nanosecond()),
	}
	return transaction, nil
}
