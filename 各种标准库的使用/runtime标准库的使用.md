

# runtime.Caller()
```go
package main

import (
	"fmt"
	"path"
	"runtime"
)

//runtime.Caller()
func rc() {
	pc, file, line, ok := runtime.Caller(1) //表示调用的层数，0 是他本身，1 是谁调用的他 2 再往上找一层
	if !ok {
		fmt.Printf("runtime.Caller() failed\n")
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	fmt.Println(funcName)        // main.main
	fmt.Println(file)            // /Users/liuqixiang/project/go_study/src/oldbody.com/day6/03runtime_demo/main.go
	fmt.Println(path.Base(file)) // main.go
	fmt.Println(line)            // 24
}

func main() {
	rc()
}
```