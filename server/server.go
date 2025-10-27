package main

import (
	proto "ChitChat/grpc"
	"context"
	"flag"
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
	flag.Parse()
	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Server STARTUP_ERROR: failed to listen on %s: %v", addr, err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterChitChatServer(grpcServer, proto.UnimplementedChitChatServer{})

	log.Printf("Server STARTUP: listening on %s", addr)

	// run server
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server ERROR: %v", err)
		}
	}()
	select {}
}
