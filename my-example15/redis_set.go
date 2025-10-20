package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/go-redis/redis/v8"
    "golang.org/x/net/context"
)

type RedisSinkConnector struct {
    consumer *kafka.Consumer
    redis    *redis.Client
    topic    string
}

func NewRedisSinkConnector(kafkaBootstrap, topic, groupID, redisAddr string) (*RedisSinkConnector, error) {
    // 連接 Kafka Consumer
    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": kafkaBootstrap,
        "group.id":          groupID,
        "auto.offset.reset": "earliest",
    })
    if err != nil {
        return nil, err
    }

    // 連接 Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: "", // 無密碼
        DB:       0,  // 使用默認 DB
    })

    return &RedisSinkConnector{
        consumer: consumer,
        redis:    redisClient,
        topic:    topic,
    }, nil
}

func (r *RedisSinkConnector) ProcessMessages() error {
    ctx := context.Background()
    
    // 訂閱主題
    err := r.consumer.SubscribeTopics([]string{r.topic}, nil)
    if err != nil {
        return err
    }

    for {
        msg, err := r.consumer.ReadMessage(-1)
        if err != nil {
            return err
        }

        var user User
        if err := json.Unmarshal(msg.Value, &user); err != nil {
            log.Printf("Error unmarshaling message: %v", err)
            continue
        }

        // 更新到 Redis
        if err := r.updateRedis(ctx, user); err != nil {
            log.Printf("Error updating Redis: %v", err)
            continue
        }

        fmt.Printf("Updated user in Redis: ID=%d, Name=%s, Money=%.2f\n", 
            user.ID, user.Name, user.Money)
    }
}

func (r *RedisSinkConnector) updateRedis(ctx context.Context, user User) error {
    // 使用 Hash 存儲用戶數據
    key := fmt.Sprintf("user:%d", user.ID)
    
    userData := map[string]interface{}{
        "id":    user.ID,
        "name":  user.Name,
        "money": user.Money,
        "updated_at": user.UpdatedAt.Format(time.RFC3339),
    }

    return r.redis.HSet(ctx, key, userData).Err()
}

func (r *RedisSinkConnector) GetUserFromRedis(ctx context.Context, userID int) (map[string]string, error) {
    key := fmt.Sprintf("user:%d", userID)
    return r.redis.HGetAll(ctx, key).Result()
}

func (r *RedisSinkConnector) Close() {
    r.consumer.Close()
    r.redis.Close()
}