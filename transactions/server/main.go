package main

import (
	"context"
	"ecommerce/codegen/transactions/v1/proto"
	v1 "ecommerce/transactions/v1"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	conn, err := pgx.Connect(context.Background(), "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
	transactionService := v1.NewTransactionService(conn)
	pb.RegisterTransactionsServer(s, transactionService)
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
