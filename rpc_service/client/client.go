package main

import (
	"context"
	"fmt"
	"io"
	"log"

	sgRPC "rpc_services/proto"

	"google.golang.org/grpc"
)

const socket string = "127.0.0.1:8021"

// const socket string = "192.168.234.128:80"

func main() {
	// grpc uses HTTP 2 which is by default uses SSL
	// we use insecure (we can also use the credentials)
	conn, err := grpc.Dial(socket, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Could not connect to : ", socket)
	}
	log.Println("Connected to ", socket)
	defer conn.Close()
	client := sgRPC.NewClientServiceClient(conn) // Using this connection to use the SimpleService
	// unary request
	makeUnaryRequest(client)

	// Client streaming
	makeClientStreaming(client)

	//Server Streaming
	makeServerStreaming(client)

	// Bi-Directional
	makeBidirectional(client)
}

func makeUnaryRequest(c sgRPC.ClientServiceClient) {
	log.Println("Making Unary Request")
	req := &sgRPC.Request{Message: "To test!"}
	log.Printf("Request - %v\n", req)
	res, err := c.UnaryRPC(context.Background(), req)
	handleAndFatalError(err)
	log.Printf("Response - %v\n", res)
}

func makeClientStreaming(c sgRPC.ClientServiceClient) {
	log.Println("Client Streaming")
	stream, err := c.ClientStreamRPC(context.Background())
	handleAndFatalError(err)

	for i := 1; i < 10; i++ {
		req := fmt.Sprintf("Request number : %d", i)
		log.Printf("Request - %v\n", req)
		stream.Send(&sgRPC.Request{Message: req})
	}
	response, err := stream.CloseAndRecv()
	handleAndFatalError(err)
	log.Printf("Response - %v\n", response)
}

func makeServerStreaming(c sgRPC.ClientServiceClient) {
	log.Println("Server Streaming")
	req := &sgRPC.Request{Message: "Need stream response"}
	log.Printf("Request - %v\n", req)
	serverStream, err := c.ServerStreamRPC(context.Background(), req)
	handleAndFatalError(err)

	for {
		response, err := serverStream.Recv()
		if err == io.EOF {
			break
		}
		handleAndFatalError(err)
		log.Printf("Response - %v\n", response)
	}
}

func makeBidirectional(c sgRPC.ClientServiceClient) {
	log.Println("Bi-Directional Streaming")
	biStream, err := c.BidirectionalRPC(context.Background())
	handleAndFatalError(err)

	// here the communication sequence is completely depends on how the server is implemented.
	// if the server is implemetend to give response to all the response at the end or
	// one after another it all compeletely depends on the implementation
	for i := 1; i < 11; i++ {
		req := fmt.Sprintf("My request %d", i)
		log.Printf("Request - %v\n", req)
		biStream.Send(&sgRPC.Request{Message: req})
		if i == 10 {
			biStream.CloseSend()
		}
		reply, err := biStream.Recv()
		handleAndPrintError(err)
		log.Printf("Response - %v\n", reply)
	}
}

func handleAndPrintError(e error) {
	if e != nil {
		log.Println(e)
	}
}

func handleAndFatalError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
