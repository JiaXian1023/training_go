kafka connect
使用 Golang 設計一個簡化的 Kafka Connect 系統來同步 MySQL 的 user 表到 Redis
系統架構
MySQL user表 → Golang Connector → Kafka → Golang Connector → Redis

CREATE TABLE user (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    money DECIMAL(10,2) DEFAULT 0.00,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_updated (updated_at)
);

-- 測試數據
INSERT INTO user (name, money) VALUES 
('Alice', 1000.50),
('Bob', 750.25),
('Charlie', 1500.75);

系統特性
實時同步: 每 2 秒檢查 MySQL 變化

容錯處理: 錯誤處理和日誌記錄

增量更新: 基於 ID 和時間戳的增量同步

數據序列化: 使用 JSON 格式傳輸

鍵值存儲: Redis 中使用 Hash 結構存儲用戶數據

這個設計提供了基本的 Kafka Connect 模式，可以根據需要擴展更多功能如：

斷點續傳

監控指標

配置管理

錯誤重試機制

