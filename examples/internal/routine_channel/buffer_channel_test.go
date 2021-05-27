package routine_channel

import (
	"log"
	"testing"
	"time"
)

// 缓冲通道在没有接受者的情况下，发送者依然可以发送指定数量的消息。不会像堵塞通道那样在发送时因为没有接受者被堵塞。
func TestBufferChannel(t *testing.T) {
	intergetCh := make(chan int, 3)
	send := func (integerCh chan int) {
		log.Print("start send")
		for i:=0;i<10;i++ {
			log.Print("about to send: ", i)
			integerCh <- i
			log.Print("sent: ", i)
			time.Sleep(time.Second)
		}
	}
	receive := func (integerCh chan int) {
		log.Print("start receive")
		time.Sleep(time.Second*3)
		for {
			v := <- integerCh
			log.Print("receive: ", v)
		}
	}
	go send(intergetCh)
	go receive(intergetCh)
	time.Sleep(time.Second*10) // 让send receive 有足够的时间去执行
}