[TOC]

# 调用系统命令打开对应资源

## 自动打开系统默认浏览器或者图片

- darwin： `open http://baidu.com`

- windows：`start http://baidu.com`
- linux： `xdg-open http://baidu.com`

```go
package main
// 打开系统默认浏览器

import (
    "fmt"
    "os/exec"
    "runtime"
)

var commands = map[string]string{
    "windows": "start",
    "darwin":  "open",
    "linux":   "xdg-open",
}

func Open(uri string) error {
    run, ok := commands[runtime.GOOS] //获取平台信息
    if !ok {
        return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
    }
    cmd := exec.Command(run, uri)
    return cmd.Start()
}

func main() {
    Open("http://baidu.com") 
    // Open("./abc.jpg")  // 打开图片
}
```



# post请求的时候head注意

```go
//post请求，必须要设定Content-Type为application/x-www-form-urlencoded，post参数才可正常传递。
req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
```

