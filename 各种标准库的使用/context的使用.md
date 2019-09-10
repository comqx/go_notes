[TOC]

前言：

> go标准库的context是一个可以主动关闭gorouting的一个模块

# 如何中断goroutine

## 使用全局变量的方式

## 使用通道的方式

```GO
package main

import (
	"fmt"
	"sync"
	"time"
)

//为啥需要context
//全局变量的方式关闭goroutine
var (
	wg       sync.WaitGroup
	notify   bool
	exitChan = make(chan bool, 1)
)

func f() {
	defer wg.Done()
	for {
		fmt.Println("lqx")
		time.Sleep(time.Millisecond * 500)
		if notify {
			break
		}
	}
}
func allVal() {
	wg.Add(1)
	go f()
	time.Sleep(time.Second * 5)
	//使用发送信号的方式，主动关闭goroutine
	notify = true
	wg.Wait()
}


//使用chan的方式关闭goroutine
func f1() {
	defer wg.Done()
FORLOOP:
	for {
		fmt.Println("lqx")
		time.Sleep(time.Millisecond * 500)
		select {
		case <-exitChan:
			break FORLOOP
		default:
		}
	}
}
func allChan() {
	wg.Add(1)
	go f1()
	time.Sleep(time.Second * 5)
	//往通道中发送中断信号，达到中断goroutine的目的
	exitChan <- true
	wg.Wait()
}

func main() {
	//allVal()
	allChan()
}
```

## 使用context方式

> ` WithCancel(parent Context) (ctx Context, cancel CancelFunc){} `

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func f2(ctx context.Context) {
FORLOOP:
	for {
		fmt.Println("context lqx is f2")
		time.Sleep(time.Millisecond * 500)
		select {
		case <-ctx.Done():
			break FORLOOP
		default:
		}
	}
}

func f1(ctx context.Context) {
	defer wg.Done()
	go f2(ctx)
FORLOOP:
	for {
		fmt.Println("context lqx is f1")
		time.Sleep(time.Millisecond * 500)
		select {
		case <-ctx.Done():
			break FORLOOP
		default:
		}
	}
}

func main() {
	//使用context，传入一个根节点，然后返回一个计数器，和一个关闭函数
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go f1(ctx)
	time.Sleep(time.Second * 5)

	//如果想关闭goroutine，那么就执行一下关闭函数
	cancel()
	wg.Wait()
}
```

# context介绍

> go1.7加入了一个新的标准库，定义了context类型，专门用来简化对于处理单个请求的多个goroutine之间的请求域的数据、取消信号、截止时间等相关操作

`context.Context`是一个接口，该接口实现了四个需要实现的方法

```go
type Context interface {
    // 需要返回当前Context被取消的时间，也就是完成工作的时间
    Deadline() (deadline time.Time, ok bool) 
    //需要返回一个channel，这个channel会在当前工作完成或者上下文被取消之后关闭，多次调用done方法会返回同一个chan
    Done() <-chan struct{} 
    //返回当前context结束的原因，只会在done返回的channel被关闭的时候才会返回非空的值
      // 如果当前Context被取消就会返回Canceled错误；
      // 如果当前Context超时就会返回DeadlineExceeded错误；
    Err() error 
    //会从context中返回键对应的值，对于同一个上下文来说，多次调用value并传入相同的key会返回相同的结果，该方法仅用于传递跨API和进程间跟请求域的数据；
    Value(key interface{}) interface{}
}
```

## 俩个默认值

### Background()和TODO()

Go内置两个函数：`Background()`和`TODO()`，这两个函数分别返回一个实现了`Context`接口的`background`和`todo`。我们代码中最开始都是以这两个内置的上下文对象作为最顶层的`partent context`，衍生出更多的子上下文对象。

- `Background()`主要用于main函数、初始化以及测试代码中，作为Context这个树结构的最顶层的Context，也就是根Context。

- `TODO()`，它目前还不知道具体的使用场景，如果我们不知道该使用什么Context的时候，可以使用这个。

`background`和`todo`本质上都是`emptyCtx`结构体类型，是一个不可取消，没有设置截止时间，没有携带任何值的Context。

## context包内的四个方法

### WithCancel

> 函数签名：`func WithCancel(parent Context) (ctx Context, cancel CancelFunc)`
>
> `WithCancel`返回带有新Done通道的父节点的副本。当调用返回的cancel函数或当关闭父上下文的Done通道时，将关闭返回上下文的Done通道，无论先发生什么情况。
>
> 取消此上下文将释放与其关联的资源，因此代码应该在此上下文中运行的操作完成后立即调用cancel

```go
//witchCancel
func gen(ctx context.Context) <-chan int {
	dst := make(chan int)
	n := 1
	go func() {
		for {
			select {
			case <-ctx.Done():
				return //// return结束该goroutine，防止泄露
			case dst <- n:
				n++
			}
		}
	}()
	return dst
}
func witchCancelMain() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}
```

### WithDeadline

> 函数签名：`func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)`
>

```go
//context.WitchDeadline
func witchDeadLineMain() {
	d := time.Now().Add(2000 * time.Millisecond) // 超时时间
	ctx, cancel := context.WithDeadline(context.Background(), d)
	// 尽管ctx会过期，但在任何情况下调用它的cancel函数都是很好的实践。
	// 如果不这样做，可能会使上下文及其父类存活的时间超过必要的时间。
	defer cancel()
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("lqx")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}
// 返回父上下文的副本，并将deadline调整为不迟于d。如果父上下文的deadline已经早于d，则WithDeadline(parent, d)在语义上等同于父上下文。当截止日过期时，当调用返回的cancel函数时，或者当父上下文的Done通道关闭时，返回上下文的Done通道将被关闭，以最先发生的情况为准。

// 取消此上下文将释放与其关联的资源，因此代码应该在此上下文中运行的操作完成后立即调用cancel。

// 定义了一个50毫秒之后过期的deadline，然后我们调用context.WithDeadline(context.Background(), d)得到一个上下文（ctx）和一个取消函数（cancel），然后使用一个select让主程序陷入等待：等待1秒后打印overslept退出或者等待ctx过期后退出。 因为ctx50秒后就过期，所以ctx.Done()会先接收到值，上面的代码会打印ctx.Err()取消原因。
```

### WithTimeout

>函数签名：`func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)`
>
>`WithTimeout`返回`WithDeadline(parent, time.Now().Add(timeout))`

```go
// context.WithTimeout

var wg sync.WaitGroup

func worker(ctx context.Context) {
LOOP:
	for {
		fmt.Println("db connecting ...")
		time.Sleep(time.Millisecond * 10) // 假设正常连接数据库耗时10毫秒
		select {
		case <-ctx.Done(): // 50毫秒后自动调用
			break LOOP
		default:
		}
	}
	fmt.Println("worker done!")
	wg.Done()
}

func withTimeOutMain() {
	// 设置一个50毫秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50) //延迟调用
	wg.Add(1)
	go worker(ctx)
	time.Sleep(time.Second * 5)
	cancel() // 通知子goroutine结束
	wg.Wait()
	fmt.Println("over")
}
```

### WithValue

> 函数签名：`func WithValue(parent Context, key, val interface{}) Context`

```go
// context.WithValue
type TraceCode string

func workers(ctx context.Context) {
	key := TraceCode("TRACE_CODE")
	traceCode, ok := ctx.Value(key).(string) // 在子goroutine中获取trace code
	if !ok {
		fmt.Println("invalid trace code")
	}
LOOP:
	for {
		fmt.Printf("worker, trace code:%s\n", traceCode)
		time.Sleep(time.Millisecond * 10) // 假设正常连接数据库耗时10毫秒
		select {
		case <-ctx.Done(): // 50毫秒后自动调用
			break LOOP
		default:
		}
	}
	fmt.Println("worker done!")
	wg.Done()
}

func withValueMain() {
	// 设置一个50毫秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	// 在系统的入口中设置trace code传递给后续启动的goroutine实现日志数据聚合
	ctx = context.WithValue(ctx, TraceCode("TRACE_CODE"), "12512312234")
	wg.Add(1)
	go workers(ctx)
	time.Sleep(time.Second * 5)
	cancel() // 通知子goroutine结束
	wg.Wait()
	fmt.Println("over")
}
```

# 使用Context的注意事项

- 推荐以参数的方式显示传递Context
- 以Context作为参数的函数方法，应该把Context作为第一个参数。
- 给一个函数方法传递Context的时候，不要传递nil，如果不知道传递什么，就使用context.TODO()
- Context的Value相关方法应该传递请求域的必要数据，不应该用于传递可选参数
- Context是线程安全的，可以放心的在多个goroutine中传递