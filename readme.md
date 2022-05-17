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

[使用 xsync.Go 运行安全的 routine](example/internal/safe_routine/main.go?embed)

## xsync.Go

使用 xsync.Go 可以避免子routine panic 后退出,并且可以通过 `as, errPanic := xsync.AsErrPanic(err)` 判断是否发生 panic

[xsync.Go 使用示例](example/internal/go/go_test.go?embed)

[xsync.Go 源码](go.go?embed)


## context

了解 context 之前我们先写一个用于测试的函数 

[代码](./example/internal/context/call_test.go?embed)

`Call` 函数 使用了 `select` 来等待 `ctx.Done()` `errCh` `resultCh` 三个 channel 的返回

### WithCancel

调用AB两个接口(http/rpc),当其中任何一个调用失败时取消调用另外一个接口

[代码](./example/internal/context/with_cancel_test.go?embed)


> 不是只有你想中途放弃，才去调用 cancel，只要你的任务正常完成了，就需要调用 cancel.
> 这样，这个 Context 才能释放它的资源（通知它的 children 处理 cancel，从它的 parent 中把自己移除.
> 甚至释放相关的 goroutine）。
> 很多人在使用这个方法的时候，都会忘记调用 cancel，切记切记，而且一定尽早释放。


> cancel 是向下传递的.
> 如果一个 WithCancel 生成的 Context 被 cancel 时，
> 如果它的子 Context（也有可能是孙，或者更低，依赖子的类型）也是 cancelCtx 类型的，
> 就会被 cancel，但是不会向上传递。
> parent Context 不会因为子 Context 被 cancel 而 cancel。

### WithTimeout

调用接口时指定最大运行时,超过时间则返回 err,这样能避免因为某些接口意外的响应慢或者网络延迟导致整个程序"卡死"

[代码](./example/internal/context/with_timeout_test.go?embed)

### WithDeadline

`context.WithDeadline` 与 `WithTimeout` 类似. `WithTimeout` 是控制多久后"触发" `<-ctx.Done()`, `WithDeadline` 是控制生米时候"触发" `<-ctx.Done()`

### WithValue

`context.WithValue()` 用于附加信息到 ctx 中, `context.Value()` 用于读取附加信息

可以在 http 的中间件中将http请求记录到某个数据库中,生成请求的ID.再将请求ID通过 `context.WithValue()` 附加到 ctx 中.
在后续遇到需要记录错误时,可以使用 `context.Value()` 查询到请求ID,便于调试.

### 实现支持 ctx 的函数

我提供了一份 [样板代码](./example/internal/context/support_ctx_test.go?embed) 供你参考


### routine 在 cancel 之后依然在执行 <a id="routineStillRunningAfterCancel"></a>

上面的示例基本上都展示了 ctx 的各种取消场景.

研究 [Call(ctx context.Context, opt Option)](./example/internal/context/call_test.go?embed) 的源码可以发现它是启动了一个新的 routine.

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

我们来看下面这段代码,并且运行它:

[代码](./example/internal/context/cancel_test.go?embed)

> ctx 的取消只是"不管"函数的运行结果,强行认定函数执行"错误",并将错误原因定义为超时.
> 即使被取消,被调用的函数的其他操作依然会继续运行


## routine 泄露 <a id="routine-leaks"></a>

> 很多原因会导致 channel 堵塞,一旦发生意料之外的持续的堵塞会导致routine一直不被释放.这种情况叫routine泄露,会导致CPU内存爆满.

主 routine退出后，系统会自动回收运行时资源，一般情况下子 routine 会自动释放
但是应该尽量避免泄露。比如在常驻服务中，比如 http server，每接收到一个请求，便会启动一次协程.
那么 子routine越来越多，每次启动的 routine 都得不到释放，内存占用和CPU会越来越高直到崩溃
避免 bug 的方法是:使用容量为1的缓冲通道
我们可通过性能分析去观测内存泄露的代码

[代码](./example/internal/routine_leaks/main.go?embed)


> 死记硬背 routine 泄露的几种情况是治标不治本,理解泄露的原因就可以通过推论得到答案.



