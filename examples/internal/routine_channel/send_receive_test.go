package routine_channel

import (
	"log"
	"testing"
	"time"
)
/*
默认情况下，通信是同步且无缓冲的：在**有**接受者接收数据之前，发送会堵塞
*/
func TestSendReceive(t *testing.T) {
	send := func (integerCh chan int) {
		log.Print("start send")
		for i:=0;i<3;i++ {
			log.Print("about to send: ", i)
			integerCh <- i
			log.Print("sent: ", i)

		}
	}
	receive := func (integerCh chan int) {
		log.Print("start receive")
		time.Sleep(time.Second*2) // 在接收方没有开始 v := <- 时，发送方的 integerCh <- i 将会堵塞
		for {
			v := <- integerCh
			log.Print("receive: ", v)
		}
	}
	intergetCh := make(chan int)
	go send(intergetCh)
	go receive(intergetCh)
	time.Sleep(time.Second*6) // 让send receive 有足够的时间去执行
}
