# http客户端与服务端

## 服务端

```go
package main

import (
	"fmt"
	"net/http"
)

func sayHello(resp http.ResponseWriter, r *http.Request) {
	var respIntContent string
	if r.Method == "POST" {
		respIntContent = "post"
	} else if r.Method == "GET" {
		respIntContent = "get"
	}
	fmt.Fprint(resp, respIntContent) //把respIntContent写入resp这个对象中，返回给http前端
}

func main() {
	http.HandleFunc("/hello", sayHello) //指定url路径以及执行的函数
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Println("start http server failed ,err", err)
	}
}
```

## 客户端

```go
func main() {
	//get 请求
	resp, err := http.Get("http://qixiang-liu.github.io")
	if err != nil {
		fmt.Println("get url failed err", err)
	}
	defer resp.Body.Close()
	//读取body
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("readall body err ", err)
	}
	fmt.Println(string(b))
}
```

