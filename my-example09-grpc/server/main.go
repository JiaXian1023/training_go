package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"github.com/yourname/gin-websocket-grpc/proto"
)

type chatServer struct {
	proto.UnimplementedChatServiceServer
}

func (s *chatServer) SendMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	log.Printf("Received message: %s from user %s", req.Content, req.UserId)
	return &proto.MessageResponse{
		Content:   fmt.Sprintf("Echo: %s", req.Content),
		UserId:    req.UserId,
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

func (s *chatServer) StreamMessages(req *proto.StreamRequest, stream proto.ChatService_StreamMessagesServer) error {
	for i := 0; i < 5; i++ {
		msg := &proto.MessageResponse{
			Content:   fmt.Sprintf("Stream message %d", i+1),
			UserId:    req.UserId,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		if err := stream.Send(msg); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterChatServiceServer(s, &chatServer{})

	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}