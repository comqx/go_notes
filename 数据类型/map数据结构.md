- [map简介](#map简介)
	- [map的声明](#map的声明)
	- [map相关操作](#map相关操作)
	- [map嵌套](#map嵌套)
	- [map是一个引用类型](#map是一个引用类型)
	- [slice of map](#slice-of-map)
	- [map排序(无序)](#map排序无序)
	- [map反转](#map反转)
# map简介

- map是引用类型，必须初始化才能使用

## map的声明

> 声明是不会分配内存的，初始化需要make
```go
var map1 map[keytype]valuetype
var a map[string]string
var a map[string]int
var a map[int]string
var a map[string]map[string]string

// map实例-1 声明一个map并初始化-1
func testMap1() {
	var a map[string]string         //声明一个map
	a = make(map[string]string, 10) //初始化，申请内存空间
	a["abc"] = "efg"
	a["abc2"] = "efg44"
	fmt.Println(a)
}

// map实例-2 声明一个map并初始化-2
func testMap2() {
	var a map[string]string = map[string]string{
		"key": "value",
	} //声明一个map,并且初始化
	a["abc"] = "efg"
	a["abc2"] = "efg44"
	fmt.Println(a)
}
```
## map相关操作

```go
  //声明一个map
	var m1 map[string]int
	fmt.Println(m1 == nil) //还没有初始化，就没有内存空间的开辟

	//初始化map
	m1 = make(map[string]int, 10) //申请map的大小

	//map增加和更新值
	m1["lqx"] = 20
	m1["刘"] = 21
	m1["lqx"] = 18
	fmt.Println(m1) //map[lqx:18 刘:21]

	//判断是否存在一个值
	fmt.Println(m1["lll"]) //如果没有值，会返回一个默认值
	v, ok := m1["lll"]     //如果不存在，ok会返回false
	fmt.Println(v, ok)

	//长度
	fmt.Println(len(m1))

	//遍历map
	for i, v := range m1 {
		fmt.Println(i, v)
	}

	// delete()函数删除键值对
	delete(m1, "lqc") //如果没有，就忽略
	delete(m1, "lqx")
	fmt.Println(m1)
```



## map嵌套

```go

// map实例-3 map嵌套
func testMap3() {
	a := make(map[string]map[string]string, 100) //声明一个map，并初始化，指定大小为100
	a["key1"] = make(map[string]string)          //初始化map里面的map
  
	a["key1"]["key2"] = "abc"
	a["key1"]["key3"] = "abc"
	fmt.Println(a)
}

// map实例-4 判断map值,查找，添加，修改
func modify(a map[string]map[string]string) {
	_, ok := a["zhangsan"]
	if !ok {
		a["zhangsan"] = make(map[string]string)
	}
	a["zhangsan"]["password"] = "123456"
	a["zhangsan"]["nickname"] = "pangpang"
	return
}
func testmap4() {
	a := make(map[string]map[string]string, 100)
	modify(a)
	fmt.Println(a)
}

// map实例－５　遍历map,删除，计算长度
func testmap5() {
	a := make(map[string]map[string]string, 100) //声明一个map，并初始化，指定大小为100
	a["key1"] = make(map[string]string)          //初始化map里面的map

	a["key1"]["key2"] = "abc"
	a["key1"]["key3"] = "abc"
	a["key1"]["key4"] = "abc"
	a["key1"]["key5"] = "abc"
	a["key2"] = make(map[string]string) //初始化map里面的map
	a["key2"]["key1"] = "abc"
	a["key2"]["key2"] = "abc"
	a["key2"]["key3"] = "abc"
	for k, v := range a {
		fmt.Println(k)
		for k1, v1 := range v {
			fmt.Println("\t", k1, v1)
		}
	}
```

## map是一个引用类型

## slice of map
```go
func testmap6() {
	// var a = make([]map[int]int, 5) //等同于下面的写法
	var a []map[int]int
	a = make([]map[int]int, 5)

	for i := 0; i < 5; i++ {
		a[i] = make(map[int]int)
	}
	a[0][10] = 100
	fmt.Println(a) //[map[10:100] map[] map[] map[] map[]]
}
```

## map排序(无序)
- 先获取所有key，把key进行排序
- 按照排序好的key,进行遍历
```go
func testmapsort() {
	// var a = make([]map[int]int, 5) //等同于下面的写法
	var a map[int]int
	a = make(map[int]int, 5)

	a[0] = 1003
	a[12] = 1012311
	a[11] = 1012123
	a[10] = 103123
	a[1] = 1012222
	a[3] = 10121
	a[5] = 10123
	fmt.Println(a) //map[0:1003 1:1012222 3:10121 5:10123 10:103123 11:1012123 12:1012311]
}
```
## map反转
- 初始化另外一个map，把key，value互换即可
```go
func testmapsort1() {
	var a map[string]int
	var b map[int]string

	a = make(map[string]int, 5)
	b = make(map[int]string, 5)
	a["abc"] = 101
	a["ccc"] = 1123123
	for k, v := range a {
		b[v] = k
	}
	fmt.Println(a, "\t", b)
}
```