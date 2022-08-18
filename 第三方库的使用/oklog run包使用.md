# oklog/run 包介绍

> oklog/run 包非常简单，只有一个类型，两个方法，共 60 行代码。其中 Group 是一组 actor，通过调用 Add 方法将 actor 添加到 Group 中。

```golang
// Package run implements an actor-runner with deterministic teardown. It is
// somewhat similar to package errgroup, except it does not require actor
// goroutines to understand context semantics. This makes it suitable for use in
// more circumstances; for example, goroutines which are handling connections
// from net.Listeners, or scanning input from a closable io.Reader.
package run

// Group collects actors (functions) and runs them concurrently.
// When one actor (function) returns, all actors are interrupted.
// The zero value of a Group is useful.
type Group struct {
	actors []actor
}

// Add an actor (function) to the group. Each actor must be pre-emptable by an
// interrupt function. That is, if interrupt is invoked, execute should return.
// Also, it must be safe to call interrupt even after execute has returned.
//
// The first actor (function) to return interrupts all running actors.
// The error is passed to the interrupt functions, and is returned by Run.
func (g *Group) Add(execute func() error, interrupt func(error)) {
	g.actors = append(g.actors, actor{execute, interrupt})
}

// Run all actors (functions) concurrently.
// When the first actor returns, all others are interrupted.
// Run only returns when all actors have exited.
// Run returns the error returned by the first exiting actor.
func (g *Group) Run() error {
	if len(g.actors) == 0 {
		return nil
	}

	// Run each actor.
	errors := make(chan error, len(g.actors))
	for _, a := range g.actors {
		go func(a actor) {
			errors <- a.execute()
		}(a)
	}

	// Wait for the first actor to stop.
	err := <-errors

	// Signal all actors to stop.
	for _, a := range g.actors {
		a.interrupt(err)
	}

	// Wait for all actors to stop.
	for i := 1; i < cap(errors); i++ {
		<-errors
	}

	// Return the original error.
	return err
}

type actor struct {
	execute   func() error
	interrupt func(error)
}
```



## 例子

> 下面例子定义了三个 actor，前两个 actor 一直等待。第三个 actor 在 3s 后结束退出。引起前两个 actor 退出。

```golang
package main

import (
	"fmt"
	"github.com/oklog/run"
	"time"
)

func main() {
	g := run.Group{}
	{
		cancel := make(chan struct{})
		g.Add(
			func() error {

				select {
				case <- cancel:
					fmt.Println("Go routine 1 is closed")
					break
				}

				return nil
			},
			func(error) {
				close(cancel)
			},
		)
	}
	{
		cancel := make(chan struct{})
		g.Add(
			func() error {

				select {
				case <- cancel:
					fmt.Println("Go routine 2 is closed")
					break
				}

				return nil
			},
			func(error) {
				close(cancel)
			},
		)
	}
	{
		g.Add(
			func() error {
				for i := 0; i <= 3; i++ {
					time.Sleep(1 * time.Second)
					fmt.Println("Go routine 3 is sleeping...")
				}
				fmt.Println("Go routine 3 is closed")
				return nil
			},
			func(error) {
				return
			},
		)
	}
	g.Run()
}

// 打印结果：
Go routine 3 is sleeping...
Go routine 3 is sleeping...
Go routine 3 is sleeping...
Go routine 3 is closed
Go routine 2 is closed
Go routine 1 is closed
```

