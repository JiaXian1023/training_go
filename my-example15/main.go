package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // 配置參數
    mysqlDSN := "user:password@tcp(localhost:3306)/testdb"
    kafkaBootstrap := "localhost:9092"
    redisAddr := "localhost:6379"
    topic := "user-updates"

    // 啟動 Source Connector
    source, err := NewMySQLSourceConnector(mysqlDSN, kafkaBootstrap, topic)
    if err != nil {
        log.Fatal("Failed to create source connector:", err)
    }
    defer source.Close()

    // 啟動 Sink Connector
    sink, err := NewRedisSinkConnector(kafkaBootstrap, topic, "redis-sink-group", redisAddr)
    if err != nil {
        log.Fatal("Failed to create sink connector:", err)
    }
    defer sink.Close()

    // 啟動處理器
    go source.Start()
    go func() {
        if err := sink.ProcessMessages(); err != nil {
            log.Fatal("Sink connector error:", err)
        }
    }()

    // 等待中斷信號
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
    <-sigchan

    log.Println("Shutting down connectors...")
}