4.1 集成Prometheus监控
go
import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	connectionsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "websocket_connections_total",
		Help: "Current number of WebSocket connections",
	})
	messagesCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "websocket_messages_total",
		Help: "Total number of WebSocket messages",
	}, []string{"direction"})
)

func init() {
	prometheus.MustRegister(connectionsGauge)
	prometheus.MustRegister(messagesCounter)
}

// 在Hub的run方法中更新指标
func (h *Hub) run() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.mu.RLock()
			connectionsGauge.Set(float64(len(h.clients)))
			h.mu.RUnlock()
		// ... 其他case
		}
	}
}

// 在main函数中添加metrics端点
http.Handle("/metrics", promhttp.Handler())

4.2 优雅关闭
go
func main() {
	// ... 初始化代码

	// 优雅关闭通道
	done := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		// 关闭所有WebSocket连接
		hub.mu.Lock()
		for client := range hub.clients {
			client.conn.WriteMessage(websocket.CloseMessage, []byte{})
			client.conn.Close()
			delete(hub.clients, client)
		}
		hub.mu.Unlock()

		// 给HTTP服务器30秒完成现有请求
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



关键优化点总结
连接管理：

使用Hub模式集中管理所有连接

读写分离，各自使用独立的goroutine

心跳机制保持连接活跃

并发控制：

使用sync.RWMutex保护共享数据

通道缓冲大小根据负载调整

限制单个连接资源使用

资源优化：

缓冲区重用(sync.Pool)

消息大小限制

连接数监控

可观测性：

Prometheus指标集成

详细日志记录

健康检查端点

稳定性：

优雅关闭处理

错误恢复机制

资源泄漏防护

这个实现可以支持数万级别的并发WebSocket连接，根据实际硬件配置和网络环境，通过调整参数可以进一步优化性能。