# 全局中间件

- 所有请求都经过此中间件

```go
package main

import (
   "github.com/gin-gonic/gin"
   "time"
   "fmt"
)

// 定义中间
func MiddleWare() gin.HandlerFunc {
   return func(c *gin.Context) {
      t := time.Now()
      fmt.Println("中间件开始执行了")
      // 设置变量到Context的key中，可以通过Get()取
      c.Set("request", "中间件")
      // 执行函数
      c.Next()
      // 中间件执行完后续的一些事情
      status := c.Writer.Status()
      fmt.Println("中间件执行完毕", status)
      t2 := time.Since(t)
      fmt.Println("time:", t2)
   }
}

func main() {
   r := gin.Default()
   // 注册中间件
   r.Use(MiddleWare())
   // {}为了代码规范
   {
      r.GET("/middleware", func(c *gin.Context) {
         // 取值
         req, _ := c.Get("request")
         fmt.Println("request:", req)
         // 页面接收
         c.JSON(200, gin.H{"request": req})
      })

   }
   r.Run(":8000")
}
```

# Next()方法源码

# 局部中间件

# 中间件练习