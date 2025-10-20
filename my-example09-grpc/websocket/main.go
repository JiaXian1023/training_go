package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/yourname/gin-websocket-grpc/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn     *websocket.Conn
	userID   string
	grpcConn *grpc.ClientConn
	client   proto.ChatServiceClient
}

func main() {
	// 创建 Gin 路由
	r := gin.Default()

	// WebSocket 路由
	r.GET("/ws/:userID", func(c *gin.Context) {
		userID := c.Param("userID")

		// 升级为 WebSocket 连接
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err)
			return
		}
		defer ws.Close()

		// 创建 gRPC 连接
		grpcConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Failed to connect to gRPC server: %v", err)
			return
		}
		defer grpcConn.Close()

		client := proto.NewChatServiceClient(grpcConn)

		// 创建客户端实例
		cli := &Client{
			conn:     ws,
			userID:   userID,
			grpcConn: grpcConn,
			client:   client,
		}

		// 处理消息
		cli.handleMessages()
	})

	// 启动 Gin 服务器
	log.Println("WebSocket gateway listening on :8080")
	r.Run(":8080")
}

func (c *Client) handleMessages() {
	// 启动 gRPC 流
	stream, err := c.client.StreamMessages(context.Background(), &proto.StreamRequest{UserId: c.userID})
	if err != nil {
		log.Printf("Failed to create stream: %v", err)
		return
	}

	// 接收 gRPC 流消息
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Stream error: %v", err)
				return
			}
			if err := c.conn.WriteJSON(msg); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}()

	// 处理 WebSocket 消息
	for {
		var msg struct {
			Content string `json:"content"`
		}
		if err := c.conn.ReadJSON(&msg); err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}

		// 通过 gRPC 发送消息
		resp, err := c.client.SendMessage(context.Background(), &proto.MessageRequest{
			Content: msg.Content,
			UserId: c.userID,
		})
		if err != nil {
			log.Printf("gRPC send error: %v", err)
			continue
		}

		// 将响应写回 WebSocket
		if err := c.conn.WriteJSON(resp); err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}
	}
}