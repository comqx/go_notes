```GO
package main

import (
	"fmt"
	"sync"
)

//channel练习
// 1. 启动一个goroutine，生成100个随机数发送到ch1
// 2. 启动一个goroutine，从ch1通道中读取值，然后计算其平方，放到ch2中
// 3. 在main中，从ch2取值打印出来

/*
知识点：
1、开启多个goroutine，执行同一个函数的时候，需要开启 sync.Once ，保证某个操作只执行一次
		var once sync.Once
		once.Do(func() { close(ch2) })
		
2、开启多个goroutine的时候，需要开启线程等待。
	var wg sync.WaitGroup
	wg.Done()
	wg.Add(3)
	wg.Wait()
	
3、只读、只写chan
*/

var wg sync.WaitGroup
var once sync.Once //确保某个操作执行一次

func f1(ch1 chan<- int) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		ch1 <- i
	}
	close(ch1)
}

func f2(ch1 <-chan int, ch2 chan<- int) {
	defer wg.Done()
	for {
		x, ok := <-ch1
		if !ok {
			break
		}
		ch2 <- x * x
	}
	once.Do(func() { close(ch2) }) // 传入的是一个匿名函数，确保某个操作只执行一次
}

func main() {
	a := make(chan int, 100)
	b := make(chan int, 100)
	wg.Add(3)
	go f1(a)
	go f2(a, b)
	go f2(a, b)
	wg.Wait()
	for ret := range b {
		fmt.Println(ret)
	}
}
```

