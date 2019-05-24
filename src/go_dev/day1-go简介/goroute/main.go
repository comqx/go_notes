package main
import (
	"time"
	
)
func main() {
	for i := 0;i < 100; i++{
		go test_goroute(i)  //go 开启多线程
	}
	time.Sleep(time.Second)
}