# goroutine之间使用chan通信

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

# 使用goroutine实现生产者-消费者模型

```go
// 使用goroutine实现生产者-消费者模型
// 1、开启一个线程池，后台挂起指定个goroutine，读取chan里面的数据，处理完成后丢到另一个chan里面
// 2、开启多个任务给chan里面输入数据
// 3、打印处理完成的chan里面的数据

/*
知识点：
1、开启多个goroutine
2、利用生产数据输入chan，goroutine消费的概念实现
*/
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("worker:%d start job:%d\n", id, j)
		time.Sleep(time.Second)
		fmt.Printf("worker:%d end job:%d\n", id, j)
		results <- j * 2
	}
}

func workerMain() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	//开启3个goroutine---消费者
	for i := 1; i <= 3; i++ {
		go worker(i, jobs, results)
	}
	//开启5个任务--生产者
	for j := 1; j <= 1000; j++ {
		jobs <- j
	}
	close(jobs)

	//输出结果---消费者
	for a := 1; a <= 5; a++ {
		fmt.Println(<-results)
	}
}
func main() {
	workerMain()
	// worker:1 start job:2
	// worker:2 start job:3
	// worker:3 start job:1
	// worker:2 end job:3
	// worker:3 end job:1
	// 6
	// 2
	// worker:2 start job:4
	// worker:3 start job:5
	// worker:1 end job:2
	// 4
	// worker:2 end job:4
	// 8
	// worker:3 end job:5
	// 10
}
```

# 实现一个计算int64随机数各位数和的程序

```go
package main
import (
	"fmt"
	"math/rand"
	"sync"
)
// 使用goroutine和channel实现一个计算int64随机数各位数和的程序。
// 1. 开启一个goroutine循环生成int64类型的随机数，发送到jobChan
// 2. 开启24个goroutine从jobChan中取出随机数计算各位数的和，将结果发送到resultChan
// 3. 主goroutine从resultChan取出结果并打印到终端输出
/*
知识点：
1、chan里面使用结构体指针类型

*/
type job struct {
	value int64
}
type result struct {
	job *job
	sum int64
}

var jobChan = make(chan *job, 100)
var resultChan = make(chan *result, 100)
var wg sync.WaitGroup

func producer(jobChan chan *job) {
	defer wg.Done()
	// 循环生成int64类型的随机数，发送到jobChan
	for {
		x := rand.Int63()
		newJob := &job{
			value: x,
		}
		jobChan <- newJob
	}
}

func consumer(jobChan <-chan *job, resultChan chan<- *result) {
	defer wg.Done()
	// 从jobChan中取出随机数计算各位数的和，将结果发送到resultChan
	for {
		job := <-jobChan
		n := job.value
		sum := int64(0)
		for n > 0 {
			sum += n % 10
			n = n / 10
		}
		newResult := &result{
			job: job,
			sum: sum,
		}
		resultChan <- newResult
	}
}

func main() {
	wg.Add(1)
	go producer(jobChan)
	// 开启24个goroutine执行保德路
	wg.Add(24)
	for i := 0; i < 24; i++ {
		go consumer(jobChan, resultChan)
	}
	// 主goroutine从resultChan取出结果并打印到终端输出
	for results := range resultChan {
		fmt.Printf("value:%d sum:%d\n", results.job.value, results.sum)
	}
	wg.Wait()
}
```



