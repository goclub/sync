package pt

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"
)

func TestPT(t *testing.T) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.SetCookie(writer, &http.Cookie{
			Name:  "count",
			Value:  "1",
		})
		_, err := writer.Write([]byte("ok")) ; if err != nil {
			log.Print(err)
			writer.WriteHeader(500) ; return
		}
	})
	http.HandleFunc("/cookie", func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(time.Second)
		cookie, err := request.Cookie("count") ; if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				// 为了演示忽略 no cookie 的情况
			} else {
				log.Print(err)
				writer.WriteHeader(500) ; return
			}
		}
		_, err = writer.Write([]byte(cookie.Value)) ; if err != nil {
			log.Print(err)
			writer.WriteHeader(500) ; return
		}
	})
	go func() {
		http.ListenAndServe(":1111", nil)
	}()
	jar, err := cookiejar.New(nil) ; if err != nil {
	    return
	}
	pt := PT{
		Steps: func() (steps []Step, err error) {
			{
				var request *http.Request
				request, err = http.NewRequest("GET", "http://127.0.0.1:1111/", nil) ; if err != nil {
				    return
				}
				steps = append(steps, Step{
					Request: request,
					Check: func(t *testing.T, step Step, resp *http.Response) {
						data, err := ioutil.ReadAll(resp.Body) ; assert.NoError(t, err)
						assert.Equal(t, string(data), "ok")
					},
				})
			}
			{
				var request *http.Request
				request, err = http.NewRequest("GET", "http://127.0.0.1:1111/cookie", nil) ; if err != nil {
					return
				}
				steps = append(steps, Step{
					Request: request,
					Check: func(t *testing.T, step Step, resp *http.Response) {
						data, err := ioutil.ReadAll(resp.Body) ; assert.NoError(t, err)
						assert.Equal(t, string(data), "1")
					},
				})
			}
			return
		},
		client:  &http.Client{
			Jar: jar,
		},
	}
	err = pt.Run(t, RunOpt{
		Times: 1000,
		Interval:  time.Millisecond * 1,
	}) ; if err != nil {
	    log.Print(err)
	}
}