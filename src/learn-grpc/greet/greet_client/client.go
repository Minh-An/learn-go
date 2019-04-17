package main

import (
	"context"
	"fmt"
	"io"
	"learn-grpc/greet/greetpb"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:1433", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Created client: %f\n", c)

	doClientStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting Unary RPC\n")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "An",
			LastName:  "Doan",
		},
	}

	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v\n", err)
	}
	log.Printf("Response from Greet %v\n", resp.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting Server Streaming RPC\n")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "An",
			LastName:  "Doan",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v\n", err)
	}
	for {
		resp, err := resStream.Recv()
		if err == io.EOF {
			break //end of stream
		}
		if err != nil {
			log.Fatalf("Error while streaming GreetManyTimes: %v\n", err)
		}
		log.Printf("Response from GreetManyTimes %v", resp.Result)
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting Client Streaming RPC\n")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{FirstName: "An"},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{FirstName: "Claire"},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{FirstName: "Megane"},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{FirstName: "George"},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{FirstName: "Brian"},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet RPC: %v\n", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending Req %v\n", req)
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Error while sending stream LongGreet %v\n", err)
		}
		time.Sleep(time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving resp LongGreet RPC: %v\n", err)
	}
	log.Printf("Response from LongGreet %v\n", resp.Result)
}
