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
			if query.Get("name") == "goclub" {
				panic("name can not be goclub")
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
