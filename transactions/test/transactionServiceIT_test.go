package test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

// Task represents a unit of work to complete. We're going to be using this in
// our example as a way to organize data that is being manipulated in
// the database.
type task struct {
	ID         uuid.UUID             `json:"id"`
	CustomerId uuid.UUID             `json:"customer_id"`
	ProductId  uuid.UUID             `json:"product_id"`
	Quantity   int                   `json:"quantity"`
	Created    timestamppb.Timestamp `json:"date_updated"`
	Amount     float64               `json:"amount"`
}

type cockroachDBContainer struct {
	testcontainers.Container
	URI string
}

//type TransactionServiceClient struct {
//	transactionServiceClient v1.TransactionService
//}

func initCockroachDB(ctx context.Context, db *sql.DB) error {
	// Actual SQL for initializing the database should probably live elsewhere
	const query = `CREATE DATABASE trx;
		CREATE TABLE if not exists trx.task  (
			 id UUID primary key,
			 customer_id UUID not null,
			 product_id UUID not null,
			 quantity int not null,
			 amount double precision not null,
			 created timestamp default current_timestamp);`
	_, err := db.ExecContext(ctx, query)
	return err
}

//func truncateCockroachDB(ctx context.Context, db *sql.DB) error {
//	const query = `TRUNCATE projectmanagement.task`
//	_, err := db.ExecContext(ctx, query)
//	return err
//}

func /*(a *TransactionServiceClient)*/ TestIntegrationDBInsertSelect(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	cdbContainer, err := startContainer(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := cdbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	db, err := sql.Open("pgx", cdbContainer.URI+"/trx")
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, initCockroachDB(ctx, db))

	// Insert data
	tsk := task{ID: uuid.New(), ProductId: uuid.New(), CustomerId: uuid.New(), Quantity: 0, Amount: 1.0}
	const insertQuery = `insert into "task" (id, customer_id, product_id, quantity, amount)
		values ($1, $2, $3, $4, $5)`
	_, err = db.ExecContext(
		ctx,
		insertQuery,
		tsk.ID,
		tsk.CustomerId,
		tsk.ProductId,
		tsk.Quantity,
		tsk.Amount,
	)
	//_, err = a.transactionServiceClient.CreateTransaction(ctx, &pb.CreateTransactionRequest{
	//	Id:         tsk.ID.String(),
	//	CustomerId: tsk.CustomerId.String(),
	//	ProductId:  tsk.ProductId.String(),
	//	Quantity:   int32(tsk.Quantity),
	//	Amount:     tsk.Amount,
	//})
	//
	//if err != nil {
	//	println("error")
	//}
	require.NoError(t, err)

	// Select data
	savedTsk := task{ID: tsk.ID}
	const findQuery = `select customer_id, product_id, quantity, amount
		from task
		where id = $1`
	row := db.QueryRowContext(ctx, findQuery, tsk.ID)
	err = row.Scan(&savedTsk.CustomerId, &savedTsk.ProductId, &savedTsk.Quantity, &savedTsk.Amount)
	require.NoError(t, err)
	assert.Equal(t, tsk.ID, savedTsk.ID)
	assert.Equal(t, tsk.CustomerId, savedTsk.CustomerId)
	assert.Equal(t, tsk.Amount, savedTsk.Amount)

}

func startContainer(ctx context.Context) (*cockroachDBContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "cockroachdb/cockroach:latest-v21.1",
		ExposedPorts: []string{"26257/tcp", "8080/tcp"},
		WaitingFor:   wait.ForHTTP("/health").WithPort("8080"),
		Cmd:          []string{"start-single-node", "--insecure"},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "26257")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("postgres://root@%s:%s", hostIP, mappedPort.Port())

	return &cockroachDBContainer{Container: container, URI: uri}, nil
}
