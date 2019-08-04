- [时间类型](#时间类型)
	- [时间格式化](#时间格式化)
	- [时间戳](#时间戳)
	- [时间加减法](#时间加减法)
	- [时间相减(sub)](#时间相减sub)
- [时间之间转换](#时间之间转换)
	- [时间戳<--->时间格式](#时间戳---时间格式)
	- [UTC 时区转换CST时区](#utc-时区转换cst时区)
- [练习题](#练习题)

# 时间类型
```go
//时间常量，以及转换
const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond //1微秒 = 1000 纳秒
	Millisecond          = 1000 * Microsecond //1毫秒 = 1000 微秒
	Second               = 1000 * Millisecond  // 1秒 = 1000 毫秒
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)
```
## 时间格式化
```go
now := time.Now()
fmt.Println(now)   //2019-08-04 11:27:49.379591 +0800 CST m=+0.000178586
fmt.Println(now.Year(), 
			now.Month(), 
			now.Day(), 
			now.Hour(), 
			now.Minute(), 
			now.Second()) //2019 August 4 11 27 49
fmt.Println(now.Date())     //2019 August 4


// 格式化的模板为Go的出生时间2006年1月2号15点04分 Mon Jan
// 24小时制
fmt.Println(now.Format("2006-01-02 15:04:05.000 Mon Jan"))
// 12小时制
fmt.Println(now.Format("2006-01-02 03:04:05.000 PM Mon Jan"))
fmt.Println(now.Format("2006/01/02 15:04"))
fmt.Println(now.Format("15:04 2006/01/02"))
fmt.Println(now.Format("2006/01/02"))

```
## 时间戳
```go
fmt.Println(now.Unix())     //1564889269  秒
fmt.Println(now.UnixNano()) //1564889269379591000 纳秒
```
## 时间加减法
```go
now = now.Add(24 * time.Hour) //时间加法, 2019 August 4 11 27 49
fmt.Println(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())  //2019 August 5 11 27 49

now = now.Add(-48 * time.Hour) //时间减法，2019 August 5 11 27 49
fmt.Println(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()) //2019 August 3 11 37 0
```
## 时间相减(sub) 
```go
now := time.Now()
//默认解析的是UTC时区
deftimeObj, err := time.Parse("2006-01-02 15:04:05", "2019-08-04 22:22:22")
if err != nil {
	fmt.Println(err)
}
fmt.Println(deftimeObj) //2019-08-04 22:22:22 +0000 UTC

//设置本地时区
loc, err := time.LoadLocation("Asia/Shanghai")
if err != nil {
	fmt.Println(err)
	return
}
//按照指定时区解析时间 //时间格式化，传入时区
timeObj, err := time.ParseInLocation("2006-01-02 15:04:05", "2019-08-04 22:22:22", loc)
if err != nil {
	fmt.Println(err)
	return
}
fmt.Println(timeObj) //2019-08-04 22:22:22 +0800 CST
//俩个时间相减(需要先把时区转换为相同的时区，再去做减法，避免加入时区时间)
td := timeObj.Sub(now)
fmt.Println(td) //22h48m1.536951s
```

# 时间之间转换
## 时间戳<--->时间格式
```go
//时间戳----->自定义时间格式
ret := time.Unix(1564845100, 0)
fmt.Println(ret)                   //2019-08-03 23:11:40 +0800 CST
fmt.Println(ret.Year(), ret.Day()) //2019 3
fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", 
			ret.Year(),
			ret.Month(), 
			ret.Day(), 
			ret.Hour(), 
			ret.Minute(), 
			ret.Second())

//自定义时间格式----->时间戳
//按照对应的格式解析字符串类型的时间
timeObj, err := time.Parse("2006-01-02 15:04:05", "2019-08-04 12:00:01")
if err != nil {
	fmt.Printf("--->%s", err)
}
fmt.Println(timeObj)        //2019-08-04 12:00:01 +0000 UTC
fmt.Println(timeObj.Unix()) //格式化时间，并转换为unix时间戳 1564920001
```
## UTC 时区转换CST时区
```go
//设置本地时区
loc, err := time.LoadLocation("Asia/Shanghai")
if err != nil {
	fmt.Println(err)
	return
}
//按照指定时区解析时间 //时间格式化，传入时区
timeObj, err := time.ParseInLocation("2006-01-02 15:04:05", "2019-08-04 22:22:22", loc)
if err != nil {
	fmt.Println(err)
	return
}
fmt.Println(timeObj) //2019-08-04 22:22:22 +0800 CST
```

# 练习题
```go
// 1、获取当前时间，格式化输出为2017/06/19 20:30:05`格式。
// 2、编写程序统计一段代码的执行耗时时间，单位精确到微秒。

func time1() {
	now := time.Now()
	fmt.Println(now.Format("2006/01/02 15:04:05")) //2019/08/04 11:53:33
}

func times(time1 func()) {
	starttime := time.Now().Nanosecond() / 1000
	time1()
	endtime := time.Now().Nanosecond() / 1000
	fmt.Println(endtime - starttime) // 199微秒

}
```