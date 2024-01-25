package main

import (
	"context"
	v1 "ecommerce/analytics/v1"
	pb "ecommerce/codegen/analytics/v1/proto"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	conn, err := pgx.Connect(context.Background(), "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
	analyticsService := v1.NewAnalyticsService(conn)
	pb.RegisterAnalyticsServer(s, analyticsService)
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
