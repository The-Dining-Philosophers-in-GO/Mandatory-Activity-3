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

	mutex       sync.Mutex                                // locking should be possible for clocking
	subscribers map[string]proto.ChitChat_SubscribeServer // clientID -> stream
	timestamp   int64
}

func (s *ChitChatServer) Subscribe(req *proto.SubscribeRequest, stream proto.ChitChat_SubscribeServer) error {
	clientID := req.GetId()
	if clientID == "" {
		return errors.New("client_id required")
	}

	s.mutex.Lock()

	if s.subscribers == nil {
		s.subscribers = make(map[string]proto.ChitChat_SubscribeServer)
	}
	s.subscribers[clientID] = stream
	s.timestamp++
	currentTime := s.timestamp

	// Create and send JOIN broadcast to all other clients
	broadcast := &proto.BroadCast{
		Type:      proto.BroadCast_JOIN,
		ClientId:  clientID,
		Timestamp: currentTime,
		Message:   "",
	}

	// Send JOIN message to ALL clients including the new one
	for id, subscriber := range s.subscribers {
		if err := subscriber.Send(broadcast); err != nil {
			log.Printf("Failed to send JOIN broadcast to %s: %v", id, err)
		}
	}

	s.mutex.Unlock()

	log.Printf("Participant %s joined Chit Chat at logical time %d", clientID, currentTime)

	<-stream.Context().Done()

	s.removeSubscriber(clientID)
	log.Printf("Participant %s disconnected at logical time %d", clientID, s.timestamp)

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
	clientID := req.GetClientId()

	s.mutex.Lock()
	_, exists := s.subscribers[clientID]
	if exists {
		delete(s.subscribers, clientID)
		s.timestamp++
	}
	currentTime := s.timestamp

	broadcast := &proto.BroadCast{
		Type:      proto.BroadCast_LEAVE,
		ClientId:  clientID,
		Timestamp: currentTime,
	}

	for _, subscriber := range s.subscribers {
		if err := subscriber.Send(broadcast); err != nil {
			println("Failed to leave:", err.Error())
		}
	}
	s.mutex.Unlock()

	log.Printf("Participant %s left Chit Chat at logical time %d", clientID, s.timestamp)

	return &proto.LeaveResponse{Ack: true}, nil
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

func (s *ChitChatServer) removeSubscriber(clientID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if subscriber exists before deleting
	if _, exists := s.subscribers[clientID]; exists {
		delete(s.subscribers, clientID)
		s.timestamp++

		// Optionally broadcast that they left unexpectedly
		broadcast := &proto.BroadCast{
			Type:      proto.BroadCast_LEAVE,
			ClientId:  clientID,
			Timestamp: s.timestamp,
		}

		for _, subscriber := range s.subscribers {
			subscriber.Send(broadcast)
		}
	}
}
