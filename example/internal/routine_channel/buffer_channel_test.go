package routine_channel

import (
	"log"
	"testing"
	"time"
)

// 缓冲通道在没有接受者的情况下，发送者依然可以发送指定数量的消息。不会像堵塞通道那样在发送时因为没有接受者被堵塞。
func TestBufferChannel(t *testing.T) {
	integerCh := make(chan int, 3)
	send := func (integerCh chan int) {
		log.Print("start send")
		for i:=1;i<7;i++ {
			log.Print("about to send: ", i)
			integerCh <- i
			log.Print("sent: ", i)
			time.Sleep(time.Second)
		}
		close(integerCh)
	}
	go send(integerCh)
	{
		log.Print("start receive")
		time.Sleep(time.Second*6)
		for {
			v, more := <-integerCh
			if more == false {
				break
			}
			log.Print("receive: ", v)
		}
	}
}