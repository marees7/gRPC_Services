package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	sgRPC "rpc_services/proto"
	"google.golang.org/grpc"
)

const socket string = ":8021"

type Server struct {
	sgRPC.UnimplementedClientServiceServer
}

func main() {
	lisn, err := net.Listen("tcp", socket)
	if err != nil {
		log.Fatalln("Errored while Listen to : ", socket, err)
	}
	log.Println("Listening at ", socket)
	s := grpc.NewServer()
	sgRPC.RegisterClientServiceServer(s, &Server{}) // registering our grpc server with our grpc service.
	err = s.Serve(lisn)
	if err != nil {
		log.Fatalln("Errored while Serving : ", socket, err)
	}
}

func (s *Server) UnaryRPC(ctx context.Context, req *sgRPC.Request) (*sgRPC.Response, error) {
	log.Println("Unary request")
	log.Printf("Request - %v\n", req)
	response := &sgRPC.Response{Message: "Here is your response"}
	log.Printf("Response - %v\n", response)
	return response, nil
}

func (s *Server) ClientStreamRPC(stream sgRPC.ClientService_ClientStreamRPCServer) error {
	log.Println("ClientStreaming RPC")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			response := &sgRPC.Response{Message: req.Message}
			log.Printf("Response - %v\n", response)
			stream.SendAndClose(response)
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Request - %v\n", req)
	}
	return nil
}

func (s *Server) ServerStreamRPC(req *sgRPC.Request, stream sgRPC.ClientService_ServerStreamRPCServer) error {
	log.Println("ServerStreaming RPC")
	log.Printf("Request- %v", req)
	for i := 1; i < 10; i++ {
		res := fmt.Sprintf("Here is the response %d", i)
		log.Printf("Response - %v", res)
		stream.Send(&sgRPC.Response{Message: res})
	}
	return nil
}

func (s *Server) BidirectionalRPC(stream sgRPC.ClientService_BidirectionalRPCServer) error {
	log.Println("StreamingBiDirectional RPC")
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Errored in stream Recv", err)
			break
		}
		log.Printf("Request - %v", msg)
		r := fmt.Sprintf("Response for your request - %v", msg.Message)
		log.Printf("Response - %v\n", r)
		stream.Send(&sgRPC.Response{Message: r})
	}
	return nil
}
