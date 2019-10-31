# 全局中间件

- 所有请求都经过此中间件

```go
package main

import (
   "github.com/gin-gonic/gin"
   "time"
   "fmt"
)
// 自定义中间件
func MiddleWare() gin.HandlerFunc{
	return func(c *gin.Context){
		t:=time.Now()
		fmt.Println("中间件开始执行了") // 1
		//设置变量到context的key中，可以通过get()取
		c.Set("request","中间件")
		//执行函数
		c.Next()
		//中间件执行完成后续的一些事情
		status :=c.Writer.Status()
		fmt.Println("中间件执行完毕",status) // 3 
		t2:=time.Since(t)// 计算时间间隔，t时间以后到现在是多长时间
		fmt.Println("time:",t2) // 4
	}
}

// 中间件
func main() {
	r:=gin.Default()
	// 注册中间件
	r.Use(MiddleWare())
	// {}为了代码规范
	{
		r.GET("/middleware",func(c *gin.Context){
			//取值
			req,_:=c.Get("request")
			fmt.Println("request:",req) // 2 
			//页面接收
			c.JSON(200,gin.H{"request":req}) // 5
		})
	}
	r.Run(":8000")
}
```
返回值：
```shell
中间件开始执行了
request: 中间件
中间件执行完毕 200
time: 62.543µs
[GIN] 2019/10/30 - 15:02:47 | 200 |      67.906µs |       127.0.0.1 | GET      /middleware
```



# Next()方法源码



# 局部中间件

```GO
		// 根路由后面是自定义的局部中间件
		r.GET("/middleware2",MiddleWare(),func(c *gin.Context){
			//取值
			req,_:=c.Get("request")
			fmt.Println("request:",req)
			//页面接收
			c.JSON(200,gin.H{"request":req})
		})
```

# 中间件练习

> 定义程序计时中间件，然后定义2个路由，执行函数后应该打印统计的执行时间，如下：

```shell
程序运行用时: 5.004030533s
[GIN] 2019/10/30 - 15:24:10 | 200 |  5.004185879s |       127.0.0.1 | GET      /shopping/index
程序运行用时: 3.003854381s
[GIN] 2019/10/30 - 15:24:30 | 200 |  3.003895339s |       127.0.0.1 | GET      /shopping/home
```

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func myTime(c *gin.Context){
	start := time.Now()
	c.Next()
	//统计时长
	since:=time.Since(start)
	fmt.Println("程序运行用时:",since)
}

// 中间件
func main() {
	r:=gin.Default()
	// 注册中间件
	r.Use(myTime)
  // 定义url组
	shoppingGroup := r.Group("shopping")
	{
		shoppingGroup.GET("/index",shopIndexHandler)
		shoppingGroup.GET("/home",shopHomeHandler)
	}
	r.Run(":8000")
}

func shopIndexHandler(c *gin.Context){
	time.Sleep(5*time.Second)
}
func shopHomeHandler(c *gin.Context){
	time.Sleep(3 * time.Second)
}
```

