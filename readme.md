---
permalink: /
sidebarBasedOnContent: true
---

# goclub/sync

> 由浅入深介绍 go 中 routine/channel/context,并提供一些便利且安全的函数

## sync.Mutex 互斥锁

使用 routine 并发操作数据时需要使用互斥锁来避免数据竞争导致的数据不一致(并发问题)

[数据竞争](./example/internal/data_race/main.go?embed)

## sync.WaitGroup 等待协同

在 [数据竞争](./example/internal/data_race/main.go?embed) 示例中就使用了 WaitGroup 来等待多个routine执行完成再退出程序

## xsync.Once 只执行一次

xsync.Once 在 sync.Once 基础之上增加了错误传递功能

[once_test](./once_test.go?embed)


## routine channel

routine channel 的理解需要大量的实践 


>  routine 是平行时空 channel 是用于传递信息的时空通道

|编号|理论|代码|
|---|---|---|
| 1 | 实现异步和并发，需要使用 go func() 开启新的 routine ,在新的 routine 中做更多的事,让主程序继续运行 | |
| 2 | routine 之间的通信，需要通过 channel 实现. channel 一定要使用 make 创建 | |
| 3 | 将 channel 在代码中"连成管道"，并且考虑发送和接收都可能出现堵塞。做到人脑都可以判断死锁。 | [发送和接收](./example/internal/routine_channel/send_receive_test.go) | 
| 4 | 理解缓冲通道带来的非堵塞特性 | [缓存通道](./example/internal/routine_channel/buffer_channel_test.go) |
| 5 | 有多个 channel 时就使用 select 防止死锁 | |
| 6 | for{} 死循环可以应用在一些会持续不断的通过 channel 发送和接收数据的场景 | |

## 不安全的 routine

[unsafe_routine](example/internal/unsafe_routine/main.go?embed)

在web服务中子 routine 如果没有通过 defer 和  recover 处理 panic 会导致整个服务中断

## 通过 defer recover 防止服务中断

[recover_routine](example/internal/recover_routine/main.go?embed)

## xsync.Go

xsync.Go 方法能提供安全的 routine, 当发生错误和panics时候可以通过 `chan xsync.ErrorRecover` 传递 error 和 recover() 

[xsync.Go 使用示例](example/internal/gofunc/go_test.go?embed)

实现很简单,感兴趣可以看看源码
[xsync.Go 源码](go.go?embed)


## context

了解 context 之前我们先写一个用于测试的函数 

[代码](./example/internal/context/call_test.go)

`Call` 函数 使用了 `select` 来等待 `ctx.Done()` `errCh` `resultCh` 三个 channel 的返回

## WithCancel

调用AB两个接口(http/rpc),当其中任何一个调用失败时取消调用另外一个接口

[代码](./example/internal/context/with_cancel_test.go)


## WithTimeout

调用接口时指定最大运行时,超过时间则返回 err,这样能避免因为某些接口意外的响应慢或者网络抖动导致整个程序"卡死"

[代码](./example/internal/context/with_timeout_test.go)

## WithDeadline

`context.WithDeadline` 与 `WithTimeout` 类似. `WithTimeout` 是控制多久后"触发" `<-ctx.Done()`, `WithDeadline` 是控制生米时候"触发" `<-ctx.Done()`

## WithValue

`context.WithValue()` 用于附加信息到 ctx 中, `context.Value()` 用于读取附加信息

可以在 http 的中间件中将http请求记录到某个数据库中,生成请求的ID.再将请求ID通过 `context.WithValue()` 附加到 ctx 中.
在后续遇到需要记录错误时,可以使用 `context.Value()` 查询到请求ID,便于调试.


## context 真的取消了其他操作吗?

上面的示例基本上都展示了 ctx 的各种取消场景.

研究 [Call(ctx context.Context, opt Option)](./example/internal/context/call_test.go) 的源码可以发现它是启动了一个新的 routine.

在新的routine中执行一些耗时操作,并且必须使用 `select` 去判断哪个通道先返回,不同的通道返回时候的处理方式不同.

```go
select {
case <- ctx.Done():
    return "", ctx.Err()
case err := <- errCh:
    return "", err
case result := <- resultCh:
    return result, nil
}
```

如果你调用的函数没有对 ctx 进行支持, ctx 主动 cancel() 和 被动 timeout 都无法让调用函数取消.**基础库基本上都支持ctx的取消操作**


我们来看下面这段代码,并且运行它:

[代码](./example/internal/context/cancel_test.go)

> ctx 的取消只是"不管"函数的运行结果,强行认定函数执行"错误",并将错误原因定义为超时.
> 即使被取消,被调用的函数的其他操作依然会继续运行
