#令牌桶限流範例
    以下是一個完整的 Go 程式碼範例，它使用 Gin 和 rate.Limiter 來限制一個 /get-ticket API 的請求頻率，在短時間內只允許 100 個請求通過。
 

#測試限流效果
 for i in {1..150}; do curl -s http://localhost:8080/get-ticket & done



#程式碼解釋
-TicketLimiter 結構體: 這個結構體包含一個 rate.Limiter 實例，是我們實現限流的核心。

-NewTicketLimiter 函式: 這是我們的建構函式，用來初始化限流器。

-rate.Limit(100): 設定每秒最多生成 100 個令牌。

-100: 設定令牌桶的容量為 100。

-重要: 這兩個參數一起決定了限流規則。如果 burst (桶容量) 比 rate (每秒生成數) 大，表示可以處理短時間內湧入的突發流量。在這個範例中，rate 和 burst 都是 100，表示在短時間內最多能處理 100 個請求，後續的請求則會被限制。



#CheckLimit 中介層:

    tl.limiter.Allow(): 這個方法會嘗試從令牌桶中取出一個令牌。

    如果成功取出，它會返回 true，程式會執行 c.Next() 繼續處理請求。

    如果令牌桶為空，它會返回 false，這時我們就返回 HTTP 429 Too Many Requests 的錯誤碼，並使用 c.Abort() 停止後續處理。

#main 函式:

    我們創建了一個 TicketLimiter 實例，並將它作為中介層應用到 /get-ticket 這個 API 路徑上。

    這樣，所有對 /get-ticket 的請求都必須先經過 CheckLimit 的檢查。


#利用 golang.org/x/time/rate 庫來實現一個高併發場景下的限流功能。