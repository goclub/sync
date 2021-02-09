# goclub/sync

> xsync

## xsync.Routine

## 不安全的 routine

[unsafe_routine](./examples/unsafe_routine/main.go)
```.go
package main

import (
	"log"
	"net/http"
	"sync"
)

func main () {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 直接使用 go func 时如果没有 defer 处理 panic 会导致整个http监听都中断。在项目中会导致服务意外中断
			// 因为没办法检查 routine 中有没有可能会导致 panic 的代码
			query := request.URL.Query()
			if query.Get("name") == "nimoc" {
				panic("name can not be nimoc")
			}
		}()
		wg.Wait()
		_, err := writer.Write([]byte("ok")) ; if err != nil {
			log.Print(err)
			writer.WriteHeader(500)
		}
	})
	addr := ":4001"
	log.Print("访问 http://127.0.0.1" + addr)
	log.Print("然后访问 http://127.0.0.1" + addr + "/?name=nimoc")
	log.Print("接着访问 http://127.0.0.1" + addr)
	log.Print("会发现第三次访问时服务已经中断了")
	log.Print(http.ListenAndServe(addr, nil))
}

```

在web服务中子 routine 如果没有通过 defer 和  recover 处理 panic 会导致整个服务中断

## 通过 defer recover 防止服务中断

[recover_routine](./examples/recover_routine/main.go)
```.go
package main

import (
	"log"
	"net/http"
	"sync"
)

func main () {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		/* 新增代码 */ var recoverValue interface{}
		go func() {
			defer wg.Done()
			/* <- 新增代码 */
			defer func() {
				r := recover()
				if r != nil {
					recoverValue = r
				}
			}()
			/* -> */
			query := request.URL.Query()
			if query.Get("name") == "nimoc" {
				panic("name can not be nimoc")
			}
		}()
		wg.Wait()
		/* <- 新增代码 */
		if recoverValue != nil {
			log.Print(recoverValue)
			writer.WriteHeader(500) ; return
		}
		/* -> */
		_, err := writer.Write([]byte("ok")) ; if err != nil {
			log.Print(err)
			writer.WriteHeader(500)
		}
	})
	addr := ":4002"
	log.Print("访问 http://127.0.0.1" + addr)
	log.Print("然后访问 http://127.0.0.1" + addr + "/?name=nimoc")
	log.Print("接着访问 http://127.0.0.1" + addr)
	log.Print("第三次访问时服务还是正常的")
	log.Print(http.ListenAndServe(addr, nil))
}

```

## 使用 xsync.Routine{}.Go() 防止服务中断 

[safe_routine](./examples/safe_routine/main.go)
```.go
package main

import (
	xsync "github.com/goclub/sync"
	"log"
	"net/http"
)

func main () {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		routine := new(xsync.Routine)
		routine.Go(func() error {
			query := request.URL.Query()
			if query.Get("name") == "nimoc" {
				panic("name can not be nimoc")
			}
			return nil
		})
		err, recoverValue := routine.Wait()
		if err != nil {
			log.Print(err)
			writer.WriteHeader(500) ; return
		}
		if recoverValue != nil {
			log.Print(recoverValue)
			writer.WriteHeader(500) ; return
		}
		_, err = writer.Write([]byte("ok")) ; if err != nil {
			log.Print(err)
			writer.WriteHeader(500); return
		}
	})
	addr := ":4003"
	log.Print("访问 http://127.0.0.1" + addr)
	log.Print("然后访问 http://127.0.0.1" + addr + "/?name=nimoc")
	log.Print("接着访问 http://127.0.0.1" + addr)
	log.Print("第三次访问时服务还是正常的")
	log.Print(http.ListenAndServe(addr, nil))
}

```

`xsync.Routine{}.Go(routine func() error)` 在 routine 前通过 defer recover 捕获了 panic ,
并通过 `xsync.Routine{}.Wait() (error, interface{})` 返回了错误和异常方便进行处理。
增加 error 的支持是因为这样能更方便的传递错误，如果没有错误的时候返回 nil 即可。