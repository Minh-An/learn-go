package main

import (
	"context"
	"fmt"
	"learn-grpc/calculator/calculatorpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:1444", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doSum(c)
}

func doSum(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Starting Unary Sum RPC\n")

	req := &calculatorpb.SumResquest{
		FirstNum:  32,
		SecondNum: 45,
	}
	resp, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v\n", err)
	}
	log.Printf("Response from Sum  %d\n", resp.Sum)
}
