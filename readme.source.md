# goclub/sync

> xsync

## 由浅入深的介绍 goroutine channel

1. 因为想要异步和并发，所以需要 goroutine
2. 因为想要通信，所以需要 channel, 记住 channel 一定要 make
3. 将 channel 在代码中连成一条线，并且考虑发送和接收都可能出现堵塞。做到人脑都可以判断死锁。
4. 理解缓冲通道带来的非堵塞特性
5. 有多个 channel 时就使用 select 防止死锁
6. for{} 死循环可以应用在一些会持续不断的通过 channel 发送和接收数据的场景


## xsync.Routine

## 不安全的 routine

[unsafe_routine|embed](examples/internal/unsafe_routine/main.go)

在web服务中子 routine 如果没有通过 defer 和  recover 处理 panic 会导致整个服务中断

## 通过 defer recover 防止服务中断

[recover_routine|embed](examples/internal/recover_routine/main.go)

## 使用 xsync.Routine{}.Go() 防止服务中断 

[safe_routine|embed](examples/internal/safe_routine/main.go)

`xsync.Routine{}.Go(routine func() error)` 在 routine 前通过 defer recover 捕获了 panic ,
并通过 `xsync.Routine{}.Wait() (error, interface{})` 返回了错误和异常方便进行处理。
增加 error 的支持是因为这样能更方便的传递错误，如果没有错误的时候返回 nil 即可。