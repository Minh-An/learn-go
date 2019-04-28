package main

import (
	"context"
	"fmt"
	"io"
	"learn-grpc/greet/greetpb"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet func was invoked w/ %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := fmt.Sprintf("Hello %s!", firstName)
	resp := &greetpb.GreetResponse{Result: result}
	return resp, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := fmt.Sprintf("Hello %s number %d", firstName, i)
		resp := &greetpb.GreetManyTimesResponse{Result: result}
		stream.Send(resp)
		time.Sleep(time.Second)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet was invoked with streaming req\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{Result: result})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream LongGreet: %v\n", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += fmt.Sprintf("Hello %s! ", firstName)
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone was invoked with streaming req\n")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error while reading client stream GreetEveryone: %v\n", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := fmt.Sprintf("Hello %s! ", firstName)
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{Result: result})
		if sendErr != nil {
			log.Printf("Error while sending client GreetEveryone: %v\n", err)
			return sendErr
		}
	}
}

func main() {
	const PORT = 1433
	listener, err := net.Listen("tcp", "0.0.0.0:1433")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	fmt.Printf("Listening on part %d...\n", PORT)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
