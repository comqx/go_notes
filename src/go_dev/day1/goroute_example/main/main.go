package main


import (
	"../goroute"
	"fmt"
)

func main() {
	var pipe chan int  //定义管道为pipe
	pipe = make(chan int,1)   //设置管道的大小
	go  goroute.Add(100,300,pipe) //开启多线程，调用goroute里面的
	sum := <- pipe 
	fmt.Println("sum=",sum)
}