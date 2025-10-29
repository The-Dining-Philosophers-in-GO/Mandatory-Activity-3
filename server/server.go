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

// server implements the gRPC service defined in our protobuff
type ChitChatServer struct {
	proto.UnimplementedChitChatServer

	mutex       sync.Mutex                                // locking should be possible for clocking
	subscribers map[string]proto.ChitChat_SubscribeServer // clientID -> stream
	timestamp   int64
}

// Subscribe handles new client connection using server-side streaming
// This method runs the entire duration of a clients connection
func (s *ChitChatServer) Subscribe(req *proto.SubscribeRequest, stream proto.ChitChat_SubscribeServer) error {
	clientID := req.GetId()
	if clientID == "" {
		return errors.New("client_id required")
	}

	s.mutex.Lock()

	if s.subscribers == nil {
		s.subscribers = make(map[string]proto.ChitChat_SubscribeServer)
	}
	//Register the clients stream for recieving broadcasts
	//Update logical clock
	s.subscribers[clientID] = stream
	s.timestamp++
	currentTime := s.timestamp

	// Create and send JOIN broadcast to all clients
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

	//WAIT HERE, until the client disconnects or is cancelled
	<-stream.Context().Done()

	//Clean up client subscribtion
	s.removeSubscriber(clientID)
	log.Printf("Participant %s disconnected at logical time %d", clientID, s.timestamp)

	return nil
}

// Publish handles chat meesages from clients
func (s *ChitChatServer) Publish(ctx context.Context, req *proto.PublishRequest) (*proto.PublishResponse, error) {
	clientID := req.GetClientId()
	message := req.GetText()

	//Message validation - reject if its over 128char long
	if len(message) > 128 {
		return nil, status.Error(codes.InvalidArgument, "Message was too long")
	}

	//Update logical clock
	s.mutex.Lock()
	s.timestamp++
	s.mutex.Unlock()

	//Create broadcast message to send to all clients
	broadcast := &proto.BroadCast{
		Type:      proto.BroadCast_CHAT,
		ClientId:  clientID,
		Message:   message,
		Timestamp: s.timestamp,
	}

	// Send message to ALL clients including the new one
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

// Leave handles client disconnections
func (s *ChitChatServer) Leave(ctx context.Context, req *proto.LeaveRequest) (*proto.LeaveResponse, error) {
	clientID := req.GetClientId()

	//Removes client from active subscriber
	s.mutex.Lock()
	_, exists := s.subscribers[clientID]
	if exists {
		delete(s.subscribers, clientID)
		s.timestamp++
	}
	currentTime := s.timestamp

	//Create leave message to send to all clients
	broadcast := &proto.BroadCast{
		Type:      proto.BroadCast_LEAVE,
		ClientId:  clientID,
		Timestamp: currentTime,
	}

	// Send leave message to ALL clients including the new one
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

	//Create TCP listener on specified port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Server STARTUP_ERROR: failed to listen on %s: %v", addr, err)
	}
	//Creates server instance
	grpcServer := grpc.NewServer()
	//Register our service implementation with the gRPC server
	proto.RegisterChitChatServer(grpcServer, &ChitChatServer{})

	log.Printf("Server STARTUP: listening on %s", addr)

	// run server
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server ERROR: %v", err)
		}
	}()
	//Block main goroutine -keep server running :)
	select {}
}

// Remove subscriber if client disconnects unexpectedly
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
