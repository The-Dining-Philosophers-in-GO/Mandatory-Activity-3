package main

import (
	proto "ChitChat/grpc"
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type ChitChatServer struct {
	proto.UnimplementedChitChatServer

	mutex       sync.Mutex // locking should be possible for clocking
	subscribers []proto.ChitChat_SubscribeServer
	timestamp   int64
}

func (s *ChitChatServer) Subscribe(req *proto.SubscribeRequest, stream proto.ChitChat_SubscribeServer) error {
	clientID := req.GetId()
	if clientID == "" {
		return errors.New("client_id required")
	}
	s.mutex.Lock()
	s.subscribers = append(s.subscribers, stream)
	s.mutex.Unlock()

	log.Printf("Participant %s joined Chit Chat at logical time %d", clientID, s.timestamp)
	s.timestamp++
	<-stream.Context().Done() // broadcasts are sent here

	return nil
}
func (s *ChitChatServer) Publish(ctx context.Context, req *proto.PublishRequest) (*proto.PublishResponse, error) {
	clientID := req.GetClientId()
	message := req.GetText()

	if len(message) > 128 {
		return nil, status.Error(codes.InvalidArgument, "Message was too long")
	}

	s.mutex.Lock()
	s.timestamp++
	s.mutex.Unlock()

	broadcast := &proto.BroadCast{
		Type:      proto.BroadCast_CHAT,
		ClientId:  clientID,
		Message:   message,
		Timestamp: s.timestamp,
	}

	s.mutex.Lock()
	for _, subscriber := range s.subscribers {
		if err := subscriber.Send(broadcast); err != nil {
			println("Failed to send to subscriber:", err.Error())
		}
	}
	s.mutex.Unlock()
	log.Printf("Server Publish received: from=%s logical_time=%d content=%q", clientID, s.timestamp, message)

	return &proto.PublishResponse{Ack: true}, nil
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
	proto.RegisterChitChatServer(grpcServer, &ChitChatServer{})

	log.Printf("Server STARTUP: listening on %s", addr)

	// run server
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server ERROR: %v", err)
		}
	}()
	select {}
}
