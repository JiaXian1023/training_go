WaitGroup的坑
Add()操作必須早於Wait()，否則會panic
Add()設置的值必須與實際等待的goroutine數量一致，否則會panic