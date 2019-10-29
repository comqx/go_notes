# gin 数据解析和绑定

## json数据解析和绑定

> 客户端传参，后端接收并解析到结构体

```go
// json数据解析和绑定
type Login struct {
	// binding 标签的意思是必须要解析，若接收值为空，则报错，必选字段
	User string `form:"username" json:"user" uri:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}
func main() {
	r:=gin.Default()

	r.POST("loginjson",func(c *gin.Context){
		//声明接收的数据结构
		var jsonData Login
		// 将request的body中数据，自动按照json格式解析到结构体
		if err := c.ShouldBindJSON(&jsonData);err !=nil{
			// 返回错误信息
			// gin.H 封装了生成json数据的工具
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		// 判断用户名密码是否正确
		if jsonData.User !="root"|| jsonData.Password !="admin"{
			c.JSON(http.StatusBadRequest,gin.H{"status":"用户名密码错误"})
			return
			}
		c.JSON(http.StatusOK,gin.H{"status":"200"})
	})
	// 3. 监听端口，默认是8080
	r.Run(":8000")
}
```

## 表单数据解析和绑定

> 通过form表单提交的数据

```go
// form 数据解析和绑定
type Login struct {
	// binding 标签的意思是必须要解析，若接收值为空，则报错，必选字段
	User string `form:"username" json:"user" uri:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}
func main() {
	r:=gin.Default()

	r.POST("formdata",func(c *gin.Context){
		//声明接收的数据结构
		var formData Login
		// Bind()默认解析并绑定form格式
		// 根据请求头中content-type自动推断
		if err := c.Bind(&formData);err !=nil{
			// 返回错误信息
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		// 判断用户名密码是否正确
		if formData.User !="root"|| formData.Password !="admin"{
			c.JSON(http.StatusBadRequest,gin.H{"status":"用户名密码错误"})
			return
			}
		c.JSON(http.StatusOK,gin.H{"status":"200"})
	})
	// 3. 监听端口，默认是8080
	r.Run(":8000")
}
<form action="http://127.0.0.1:8000/formdata" method="post" enctype="application/form-data">
    用户名：<input type="text" name="username" >
    密&nbsp&nbsp码: <input type="password" name="password">
    <input type="submit" value="提交">
</form>
```

## uri数据解析和绑定

> 通过uri传参解析数据

```go
type Login struct {
	// binding 标签的意思是必须要解析，若接收值为空，则报错，必选字段
	User string `form:"username" json:"user" uri:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}
func main() {
	r:=gin.Default()

	r.GET("/:user/:password",func(c *gin.Context){
		//声明接收的数据结构
		var uridata Login
		// Bind()默认解析并绑定form格式
		// 根据请求头中content-type自动推断
		if err := c.ShouldBindUri(&uridata);err !=nil{
			// 返回错误信息
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		// 判断用户名密码是否正确
		if uridata.User !="root"|| uridata.Password !="admin"{
			c.JSON(http.StatusBadRequest,gin.H{"status":"用户名密码错误"})
			return
		}
		c.JSON(http.StatusOK,gin.H{"status":"200"})
	})
	// 3. 监听端口，默认是8080
	r.Run(":8000")
}
// curl http://127.0.0.1:8000/root/admin
// {"status":"200"}
```

