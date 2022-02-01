package routine_channel

import (
	"log"
	"testing"
	"time"
)
/*
默认情况下，通信是同步且无缓冲的：在有接受者接收数据之前，发送会堵塞
*/
func TestSendReceive(t *testing.T) {
	send := func (integerCh chan int) {
		log.Print("start send")
		for i:=1;i<4;i++ {
			log.Print("about to send: ", i)
			integerCh <- i
			log.Print("sent: ", i)
		}
		// 发送三次后关闭通道
		close(integerCh)
	}
	integerCh := make(chan int)
	go send(integerCh)
	{

		log.Print("start receive")
		time.Sleep(time.Second*2) // 在接收方没有开始 v := <- 时，发送方的 integerCh <- i 将会堵塞
		for {
			v, more := <- integerCh
			// 如果通道关闭了没有更多的信息传递则退出循环等待
			if more == false {
				break
			}
			log.Print("receive: ", v)
		}
	}
}
