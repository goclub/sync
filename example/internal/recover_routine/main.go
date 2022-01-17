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

		/* <- 新增的代码 */
		var recoverValue interface{}
		/* ----------> */

		go func() {
			defer wg.Done()

			/* <- 新增的代码 */
			defer func() {
				r := recover()
				if r != nil {
					recoverValue = r
				}
			}()
			/* ----------> */

			query := request.URL.Query()
			if query.Get("name") == "nimoc" {
				panic("name can not be nimoc")
			}
		}()
		wg.Wait()

		/* <- 新增的代码 */
		if recoverValue != nil {
			log.Print(recoverValue)
			writer.WriteHeader(500) ; return
		}
		/* ----------> */

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