# gin和django对比

## 中间件

利用**函数调用栈**后进先出的特点，巧妙的完成中间件在自定义处理函数完成的后处理的操作。

django它的处理方式是定义个类，请求处理前的处理的定义一个方法，请求处理后的处理定义一个方法。

gin的方式更灵活，但django的方式更加清晰。



## 请求参数绑定

对于获取请求内容，在模型绑定当中，有以下的场景

- 绑定失败是用户自己处理还是框架统一进行处理
- 用户需是否需要关心请求的内容选择不同的绑定器

在gin框架的对于这些场景给出的答案是：提供不同的方法，满足以上的需求。这里的关键点还是在于使用场景是怎样的。

```GO
// 自动更加请求头选择不同的绑定器对象进行处理
func (c *Context) Bind(obj interface{}) error {
    b := binding.Default(c.Request.Method, c.ContentType())
    return c.MustBindWith(obj, b)
}

// 绑定失败后，框架会进行统一的处理
func (c *Context) MustBindWith(obj interface{}, b binding.Binding) (err error) {
    if err = c.ShouldBindWith(obj, b); err != nil {
        c.AbortWithError(400, err).SetType(ErrorTypeBind)
    }

    return
}

// 用户可以自行选择绑定器，自行对出错处理。自行选择绑定器，这也意味着用户可以自己实现绑定器。
// 例如：嫌弃默认的json处理是用官方的json处理包，嫌弃它慢，可以自己实现Binding接口
func (c *Context) ShouldBindWith(obj interface{}, b binding.Binding) error {
    return b.Bind(c.Request, obj)
}
```

