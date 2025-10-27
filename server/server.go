package main

import (
	proto "ChitChat/grpc"
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type ChitChatServer struct {
	proto.UnimplementedChitChatServer
}

func (s *ChitChatServer) Subscribe(req *proto.SubscribeRequest, stream proto.ChitChat_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (s *ChitChatServer) Publish(ctx context.Context, req *proto.PublishRequest) (*proto.PublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (s *ChitChatServer) Leave(ctx context.Context, req *proto.LeaveRequest) (*proto.LeaveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}

func main() {
	server := &ChitChatServer{}
	//server.start_server()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterChitChatServer(grpcServer, server)
	grpcServer.Serve(lis)
}

func (s *ChitChatServer) start_server() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Did not work")
	}

	proto.RegisterChitChatServer(grpcServer, s)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work")
	}

}
