# goclub/sync

> xsync

## sync.Mutex 互斥锁

使用 routine 并发操作数据时需要使用互斥锁来避免数据竞争导致的数据不一致(并发问题)

[数据竞争](./example/internal/data_race/main.go?embed)

## sync.WaitGroup 等待协同

在 [数据竞争](./example/internal/data_race/main.go?embed) 示例中就使用了 WaitGroup 来等待多个routine执行完成再退出程序

## xsync.Once 只执行一次

xsync.Once 在 sync.Once 基础之上增加了错误传递功能

[once_test](./once_test.go?embed)

## xsync.Routine

## 不安全的 routine

[unsafe_routine](example/internal/unsafe_routine/main.go?embed)

在web服务中子 routine 如果没有通过 defer 和  recover 处理 panic 会导致整个服务中断

## 通过 defer recover 防止服务中断

[recover_routine](example/internal/recover_routine/main.go?embed)

## 使用 xsync.Routine{}.Go() 防止服务中断

[safe_routine](example/internal/safe_routine/main.go?embed)

`xsync.Routine{}.Go(routine func() error)` 在 routine 前通过 defer recover 捕获了 panic ,
并通过 `xsync.Routine{}.Wait() (error, interface{})` 返回了错误和异常方便进行处理。
增加 error 的支持是因为这样能更方便的传递错误，如果没有错误的时候返回 nil 即可。

## routine channel

1. 因为想要异步和并发，所以需要 goroutine
2. 因为想要通信，所以需要 channel, 记住 channel 一定要 make
3. 将 channel 在代码中连成一条线，并且考虑发送和接收都可能出现堵塞。做到人脑都可以判断死锁。
4. 理解缓冲通道带来的非堵塞特性
5. 有多个 channel 时就使用 select 防止死锁
6. for{} 死循环可以应用在一些会持续不断的通过 channel 发送和接收数据的场景
