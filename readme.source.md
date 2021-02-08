# goclub/sync

> xsync

## xsync.Routine

## 不安全的 routine

[unsafe_routine|embed](./examples/unsafe_routine/main.go)

在web服务中子 routine 如果没有通过 defer 和  recover 处理 panic 会导致整个服务中断

## 通过 defer recover 防止服务中断

[recover_routine|embed](./examples/recover_routine/main.go)

## 使用 xsync.Routine{}.Go() 防止服务中断 

[safe_routine|embed](./examples/safe_routine/main.go)

`xsync.Routine{}.Go(routine func() error)` 在 routine 前通过 defer recover 捕获了 panic ,
并通过 `xsync.Routine{}.Wait() (error, interface{})` 返回了错误和异常方便进行处理。
增加 error 的支持是因为这样能更方便的传递错误，如果没有错误的时候返回 nil 即可。