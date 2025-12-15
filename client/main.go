package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/mvk12/grpc-demo/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Unary RPC
func GetTestCall(c pb.TestServiceClient) {
	req := &pb.GetTestRequest{
		Id: 1,
	}
	res, err := c.GetTest(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response from UnaryCall: ID=%d, Title=%s\n", res.GetId(), res.GetTitle())
}

// Client Streaming RPC
func CreateQuestionsCall(c pb.TestServiceClient) {
	stream, err := c.CreateQuestions(context.Background())
	if err != nil {
		panic(err)
	}

	questions := []*pb.Question{
		{TestId: 1, Question: "What is gRPC?", Answer: "None"},
		{TestId: 1, Question: "Explain Protocol Buffers.", Answer: "None"},
		{TestId: 1, Question: "What are the benefits of using gRPC?", Answer: "None"},
	}

	for _, question := range questions {
		if err := stream.Send(question); err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server Response: %v\n", res.Ok)
}

// Server Streaming RPC
func GetStudentsPerTestCall(c pb.TestServiceClient) {
	req := &pb.GetStudentsPerTestRequest{
		TestId: 1,
	}
	stream, err := c.GetStudentsPerTest(context.Background(), req)
	if err != nil {
		panic(err)
	}

	fmt.Println("Students enrolled in Test ID 1:")
	for {
		student, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		fmt.Printf("ID=%d, Name=%s\n", student.GetId(), student.GetName())
	}
}

// Bidirectional Streaming RPC
func TakeTestCall(c pb.TestServiceClient) {
	md := metadata.New(map[string]string{
		"x-test-id": "1",
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := c.TakeTest(ctx)
	if err != nil {
		panic(err)
	}

	waitChannel := make(chan struct{})

	go func() {
		for {
			question, err := stream.Recv()
			if err == io.EOF {
				log.Println("Test completed.")
				break
			}

			if err != nil {
				log.Fatalf("Error receiving question: %v", err)
				break
			}

			fmt.Printf("Question ID=%d: %s\n", question.GetId(), question.GetQuestion())

			answer := &pb.TakeTestRequest{
				Answer: "Sample Answer",
			}

			if err := stream.Send(answer); err != nil {
				log.Fatalf("Error sending answer: %v", err)
				break
			}

			time.Sleep(1 * time.Second)
		}

		close(waitChannel)
	}()
	<-waitChannel
}

func main() {
	cc, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer cc.Close()

	c := pb.NewTestServiceClient(cc)

	GetTestCall(c)
	CreateQuestionsCall(c)
	GetStudentsPerTestCall(c)
	TakeTestCall(c)
}
