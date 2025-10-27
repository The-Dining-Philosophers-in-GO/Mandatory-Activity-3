package main

import (
	proto "ChitChat/grpc"
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func main() {
	flag.Parse()
	clientID := "Bob"
	conn, err := grpc.NewClient(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}

	client := proto.NewChitChatClient(conn)
	subreq := &proto.SubscribeRequest{}

	go func() {
		fmt.Println("Inside client go func")
		stream, _ := client.Subscribe(context.Background(), subreq)
		for {
			broadcast, err := stream.Recv() // Receive broadcasts
			if err != nil {
				return
			}
			fmt.Printf("[%d] %s\n", broadcast.Timestamp, broadcast.Message)
		}
	}()

	stdin := bufio.NewScanner(os.Stdin)
	fmt.Println("Type messages and press Enter to publish. Type '/leave' to exit.")
	for stdin.Scan() {
		line := stdin.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "/leave" {
			// call Leave RPC then exit
			_, err := client.Leave(context.Background(), &proto.LeaveRequest{ClientId: clientID})
			if err != nil {
				log.Printf("Client LEAVE_RPC_ERROR: %v", err)
			}
			log.Printf("Client SHUTDOWN: id=%s initiated leave", clientID)
			// Sleep briefly to allow leave broadcast to flow and then exit
			time.Sleep(200 * time.Millisecond)
			return
		}

		// Publish
		response, err := client.Publish(context.Background(), &proto.PublishRequest{ClientId: clientID, Text: line})
		if err != nil {
			log.Printf("Client PUBLISH_ERROR: %v", err)
			continue
		}
		if !response.Ack {
			log.Printf("Client PUBLISH_REJECTED: reason=%s", response.Error)
			continue
		}
		log.Printf("Client PUBLISH_SENT: id=%s content=%q", clientID, line)
	}

	if stdin.Err() != nil {
		log.Printf("Client STDIN_ERROR: %v", stdin.Err())
	}
}
