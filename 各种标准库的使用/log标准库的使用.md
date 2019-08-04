

# log库的使用
```go
package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

//runtime.Caller()，处理err问题
func rc(err error) {
	if err != nil {
		pc, file, line, ok := runtime.Caller(1) //表示调用的层数，0 是他本身，1 是谁调用的他 2 再往上找一层
		if !ok {
			fmt.Printf("runtime.Caller() failed\n")
			return
		}
		funcName := runtime.FuncForPC(pc).Name()
		fmt.Printf("funcName:%s, file:%s, line:%d, cont:%s", funcName, path.Base(file), line, err)
	}
}

func main() {
	fileObj, err := os.OpenFile("./xx.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	rc(err)              //funcName: main.main,file:main.go, line:25,cont:open ./xx.log: no such file or directory<nil>
	fmt.Println(fileObj) //文件句柄，也就是文件的指针
	log.SetOutput(fileObj)
	for {
		log.Println("这是一条测试的日志")
		time.Sleep(time.Second * 3)
	}
}
```