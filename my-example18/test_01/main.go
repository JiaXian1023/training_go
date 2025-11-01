package main
import(
	"fmt"
	"sync"
	"sync/atomic"
)
func main() {
    var counter int32
    var wg sync.WaitGroup
    
    // 启动 100 个 goroutine 同时增加计数器
    for i := 0; i < 100; i++ {
        wg.Add(1)
		fmt.Println("i",i)
        go func() {
			fmt.Println("2i",i)
            defer wg.Done()
			//fmt.Println("i",i)
            for j := 0; j < 1000; j++ {
                atomic.AddInt32(&counter, 1)
            }
        }()
    }
    
    wg.Wait()
    fmt.Printf("最终计数器值: %d\n", counter) // 总是 100000
}