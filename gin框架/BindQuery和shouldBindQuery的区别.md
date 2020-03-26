# BindQuery
> BindQuery在请求过程中， 如果参数错误会直接抛异常 返回400状态
```go
// BindQuery is a shortcut for c.MustBindWith(obj, binding.Query).
func (c *Context) BindQuery(obj interface{}) error {
	return c.MustBindWith(obj, binding.Query)
}

// MustBindWith binds the passed struct pointer using the specified binding engine.
// It will abort the request with HTTP 400 if any error ocurrs.
// See the binding package.
func (c *Context) MustBindWith(obj interface{}, b binding.Binding) (err error) {
	if err = c.ShouldBindWith(obj, b); err != nil {
		c.AbortWithError(http.StatusBadRequest, err).SetType(ErrorTypeBind)
	}
	return
}
```

# ShouldBindQuery

> ShouldBindQuery 在请求过程中，对参数检测不做处理
```go
// ShouldBindQuery is a shortcut for c.ShouldBindWith(obj, binding.Query).
func (c *Context) ShouldBindQuery(obj interface{}) error {
	return c.ShouldBindWith(obj, binding.Query)
}
// ShouldBindWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func (c *Context) ShouldBindWith(obj interface{}, b binding.Binding) error {
	return b.Bind(c.Request, obj)
}
```