package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"learn-grpc/calculator/calculatorpb"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumResquest) (*calculatorpb.SumResponse, error) {
	resp := &calculatorpb.SumResponse{
		Sum: req.GetFirstNum() + req.GetSecondNum(),
	}
	return resp, nil
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:1444")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	fmt.Printf("Listening on part %d...\n", 1444)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
