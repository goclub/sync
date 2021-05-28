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