package gameshared

import (
	"fmt"
	"sanguosha.com/baselib/ioservice"
	"testing"
	"time"
)

func Test_IOService(t *testing.T) {
	worker := ioservice.NewIOService("hello", 10240)
	worker.Init()
	worker.Run()

	worker.Post(func() {
		fmt.Println("woker start")
		time.Sleep(time.Second * 20)
		fmt.Println("woker end")
	})

	fmt.Println("app start")
	worker.Fini()
	fmt.Println("app end")
}
