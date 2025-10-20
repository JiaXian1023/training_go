package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 配置常量
const (
	MaxConnections      = 100000
	MessageBufferSize   = 256
	BroadcastChannelSize = 1024
	WriteWait           = 10 * time.Second
	PongWait           = 60 * time.Second
	PingPeriod         = (PongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		// 生产环境应验证来源
		return true
	},
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	info ClientInfo
}

type ClientInfo struct {
	ID        string
	IP        string
	UserAgent string
	Connected time.Time
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
	shutdown   chan struct{}
}

var (
	connectionsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "websocket_connections_total",
		Help: "Current number of WebSocket connections",
	})
	messagesCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "websocket_messages_total",
		Help: "Total number of WebSocket messages",
	}, []string{"direction"})
	connectionErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "websocket_connection_errors_total",
		Help: "Total number of WebSocket connection errors",
	}, []string{"type"})
)

func init() {
	prometheus.MustRegister(connectionsGauge)
	prometheus.MustRegister(messagesCounter)
	prometheus.MustRegister(connectionErrors)
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, BroadcastChannelSize),
		shutdown:   make(chan struct{}),
	}
}

func (h *Hub) run() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected: %s", client.info.ID)
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected: %s", client.info.ID)
		case message := <-h.broadcast:
			h.mu.RLock()
			clients := make([]*Client, 0, len(h.clients))
			for client := range h.clients {
				clients = append(clients, client)
			}
			h.mu.RUnlock()

			for _, client := range clients {
				select {
				case client.send <- message:
					messagesCounter.WithLabelValues("outbound").Inc()
				default:
					h.unregister <- client
				}
			}
		case <-ticker.C:
			h.mu.RLock()
			connectionsGauge.Set(float64(len(h.clients)))
			h.mu.RUnlock()
		case <-h.shutdown:
			h.mu.Lock()
			for client := range h.clients {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				client.conn.Close()
				close(client.send)
				delete(h.clients, client)
			}
			h.mu.Unlock()
			return
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		connectionErrors.WithLabelValues("upgrade").Inc()
		log.Println(err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, MessageBufferSize),
		info: ClientInfo{
			ID:        generateClientID(),
			IP:        r.RemoteAddr,
			UserAgent: r.UserAgent(),
			Connected: time.Now(),
		},
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512) // 限制消息大小
	c.conn.SetReadDeadline(time.Now().Add(PongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				connectionErrors.WithLabelValues("read").Inc()
				log.Printf("Read error: %v", err)
			}
			break
		}
		messagesCounter.WithLabelValues("inbound").Inc()
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func generateClientID() string {
	// 实现一个高效的ID生成器
	return "client-" + time.Now().Format("20060102-150405.000")
}

func main() {
	hub := newHub()
	go hub.run()

	// 注册HTTP路由
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 配置HTTP服务器
	server := &http.Server{
		Addr:              ":8098",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	// 优雅关闭
	done := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		// 通知hub关闭所有连接
		close(hub.shutdown)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown server: %v", err)
		}
		close(done)
	}()

	log.Println("WebSocket server started on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not start server: %v", err)
	}

	<-done
	log.Println("Server stopped")
}