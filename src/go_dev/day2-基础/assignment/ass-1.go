package main

import (
	"fmt"
	"time"
)

const (
	Man    = 1
	Female = 2
)

func main() {
	for {
		second := time.Now().Unix() //获取当前时间
		if second%Female == 0 {   //当前时间取余
			fmt.Println("female") 
		} else {
			fmt.Println("man")
		}
		time.Sleep(1000 * time.Microsecond)  //休息1秒
	}
}
