package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/confluentinc/confluent-kafka-go/kafka"
    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Money float64 `json:"money"`
    UpdatedAt time.Time `json:"updated_at"`
}

type MySQLSourceConnector struct {
    db        *sql.DB
    producer  *kafka.Producer
    topic     string
    lastID    int
}

func NewMySQLSourceConnector(mysqlDSN, kafkaBootstrap, topic string) (*MySQLSourceConnector, error) {
    // 連接 MySQL
    db, err := sql.Open("mysql", mysqlDSN)
    if err != nil {
        return nil, err
    }

    // 連接 Kafka
    producer, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": kafkaBootstrap,
    })
    if err != nil {
        return nil, err
    }

    return &MySQLSourceConnector{
        db:       db,
        producer: producer,
        topic:    topic,
        lastID:   0,
    }, nil
}

func (m *MySQLSourceConnector) PollChanges() error {
    query := `
        SELECT id, name, money, updated_at 
        FROM user 
        WHERE id > ? OR updated_at > DATE_SUB(NOW(), INTERVAL 1 MINUTE)
        ORDER BY id
    `
    
    rows, err := m.db.Query(query, m.lastID)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Money, &user.UpdatedAt)
        if err != nil {
            return err
        }

        // 發送到 Kafka
        if err := m.sendToKafka(user); err != nil {
            return err
        }

        // 更新最後處理的 ID
        if user.ID > m.lastID {
            m.lastID = user.ID
        }
    }

    return nil
}

func (m *MySQLSourceConnector) sendToKafka(user User) error {
    data, err := json.Marshal(user)
    if err != nil {
        return err
    }

    message := &kafka.Message{
        TopicPartition: kafka.TopicPartition{
            Topic:     &m.topic,
            Partition: kafka.PartitionAny,
        },
        Value: data,
        Key:   []byte(fmt.Sprintf("user_%d", user.ID)),
    }

    return m.producer.Produce(message, nil)
}

func (m *MySQLSourceConnector) Start() {
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        if err := m.PollChanges(); err != nil {
            log.Printf("Error polling changes: %v", err)
        }
    }
}

func (m *MySQLSourceConnector) Close() {
    m.db.Close()
    m.producer.Close()
}