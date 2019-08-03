- [struct](#struct)
	- [struct的声明](#struct的声明)
	- [struct的定义](#struct的定义)
	- [struct的初始化](#struct的初始化)
		- [先声明再赋值](#先声明再赋值)
		- [声明同时初始化](#声明同时初始化)
			- [键值对初始化](#键值对初始化)
			- [值列表初始化](#值列表初始化)
			- [注意事项](#注意事项)
	- [匿名结构体](#匿名结构体)
	- [指针类型结构体初始化](#指针类型结构体初始化)
		- [先声明再赋值](#先声明再赋值-1)
		- [声明同时初始化](#声明同时初始化-1)
	- [结构体是值类型](#结构体是值类型)
	- [结构体内存布局](#结构体内存布局)
	- [使用构造函数初始化结构体](#使用构造函数初始化结构体)
		- [方法和接收者(值和指针)](#方法和接收者值和指针)
		- [值接收者和指针接收者的区别](#值接收者和指针接收者的区别)
	- [给自定义类型添加方法](#给自定义类型添加方法)
	- [嵌套结构体](#嵌套结构体)
		- [嵌套匿名结构体](#嵌套匿名结构体)
		- [嵌套结构体字段名冲突](#嵌套结构体字段名冲突)
	- [结构体的继承](#结构体的继承)
	- [结构体字段的可见性](#结构体字段的可见性)
	- [结构体与json序列化](#结构体与json序列化)
		- [结构体的标签（tag）](#结构体的标签tag)
		- [序列化注意事项](#序列化注意事项)
	- [链表的定义](#链表的定义)
		- [链表尾部插入法](#链表尾部插入法)
# struct

- 用来自定义复杂数据结构
- struct里面可以包含多个字段(属性)
- struct类型可以定义方法，注意和函数的区分
- struct类型是值类型
- struct类型可以嵌套
- go语言没有class类型，只有struct类型

## struct的声明

> 使用`type`和`struct`关键字来定义结构体

```go
type 类型名 struct{ //类型名：标识自定义结构体的名称，在同一个包内不能重复
  字段名 字段类型 //字段名： 表示结构体字段名，结构体中的字段名必须唯一
  字段名 字段类型 //字段类型：表示结构体字段的具体类型
}
```

## struct的定义

> Struct 中字段访问，和其他语言一样，使用点

```go
type 结构体名 struct{
  字段1 字段1的类型
  字段2 字段2的类型
}

type person struct{
  name，city string
  age int8
}
结构体用来描述一组值，比如一个人有名字，年龄，居住城市等
```

## struct的初始化

### 先声明再赋值

```go
//定义struct
type person struct {
	name   string
	age    int
	gender string
	hobby  []string
}

//先声明，再赋值,没有赋值则使用零值
var p person // 声明一个person类型的变量p
p.name = "元帅" //选择字段进行赋值
p.age = 18
fmt.Println(p)
```

### 声明同时初始化

#### 键值对初始化

```go
// 键值对初始化
	var p2 = person{
		name: "冠华",
		age:  15,
	}
	fmt.Printf("p2=%#v\n", p2) //{p2=main.person{name:"冠华", age:15, gender:"", hobby:[]string(nil)}
```

#### 值列表初始化

```go
// 值列表初始化，必须结构体的所有字段
var p3 = person{
	"理想",
	100,
}
fmt.Println(p3)
```

#### 注意事项

1. 键值对初始化和值列表初始化，两者不能混用
2. 没有赋值的字段会使用对应类型的零值.

## 匿名结构体

```go
//匿名结构体。定义，并声明结构体，并初始化
func s1(){
  var s struct {
    x string
    y int
  }{10,20}
  
  fmt.Printf("type:%T, value:%v\n",s,s) 
}
```
## 指针类型结构体初始化

- 结构体是值类型,赋值的时候都是拷贝.

- 当结构体字段较多的时候,为了减少内存消耗可以传递结构体指针.

### 先声明再赋值

```GO
//定义struct
type person struct {
	name   string
	age    int
	gender string
	hobby  []string
}

//指针类型赋值
var p2 = new(person) //使用new 申请一块内存，返回的是一个指针
(*p2).name = "leix"
p2.gender="nan" //语法糖，可以这么写

fmt.Println(p2) //&{理想 0 保密 []}
fmt.Printf("%T\n", p2)  //p2是一个 *main.person类型的指针
fmt.Printf("%p\n", p2)  //0xc00007c080 这个是p2  保存 的内存地址
fmt.Printf("%p\n", &p2) //0xc00009e000  这个是p2指针的内存地址

```

### 声明同时初始化

```GO
//使用键值对的形式对指针类型struct初始化,如果对某个字段没有赋值，那么就是零值
	p3 := &person{
		name: "冠华",
		age:  15,
	}
	fmt.Printf("p3=%#v\n", p3) //p3=&main.person{name:"冠华", age:15, gender:"", hobby:[]string(nil)}

// 使用值列表的形式对指针类型struct初始化, 值的顺序要和结构体定义时字段的顺序一致
	p4 := &person{
		"小王子",
		"男",
	}
	fmt.Printf("%#v\n", p4)
```

## 结构体是值类型

```go
//结构体
type person struct {
	name   string
	age    int
	gender string
	hobby  []string
}

//go语言中函数传参永远传的是拷贝，ctrl+c ctrl+v
func f(x person) {
	x.name = "Liuqixiang"
}
func f1(x *person) {
	// (*x).gender = "nv"   //根据内存地址找到那个原变量,修改的就是原来的变量
	x.name = "指针liuqixiang" //go 语言中的语法糖，这里可以省略(*x)为x
}
func main() {
	p := person{
		name: "lqx",
		age:  20,
	}
	f(p)
	fmt.Println(p.name) //lqx
	f1(&p)
	fmt.Println(p.name) //指针liuqixiang
}
```

## 结构体内存布局

> 结构体占用一块连续的内存空间

```go
type test struct {
	a int8
	b int8
	c int8
	d int8
}
n := test{
	1, 2, 3, 4,
}
fmt.Printf("n.a %p\n", &n.a)
fmt.Printf("n.b %p\n", &n.b)
fmt.Printf("n.c %p\n", &n.c)
fmt.Printf("n.d %p\n", &n.d)
//输出：
n.a 0xc0000a0060
n.b 0xc0000a0061
n.c 0xc0000a0062
n.d 0xc0000a0063
```

【进阶知识点】关于Go语言中的内存对齐推荐阅读:[在 Go 中恰到好处的内存对齐](https://segmentfault.com/a/1190000017527311?utm_campaign=studygolang.com&utm_medium=studygolang.com&utm_source=studygolang.com)

## 使用构造函数初始化结构体

- 构造函数：约定成俗用new开头
- 返回的是结构体还是结构体指针
- 当结构体比较大的时候尽量使用结构体指针，减少程序的内存开销

```GO
//定义俩个结构体，person/dog
type person struct {
	name string
	age  int
}
type dog struct {
	name string
}
//定义一个构造函数，返回的是一个person类型的指针
//返回一个结构体变量的函数,为了实例化结构体的时候更省事儿.
func newPerson(name string, age int) *person {
	return &person{
		name: name,
		age:  age,
	}
}
func newDog(name string) *dog {
	return &dog{
		name: name,
	}
}
func main() {
	p1 := newPerson("lqx", 18)
	p2 := newPerson("abc", 20)
	fmt.Println(p1, p2) //&{lqx 18} &{abc 20}

	d1 := newDog("wangdog")
	fmt.Println(d1) //&{wangdog}
}

```

### 方法和接收者(值和指针)

```go
//方法的定义
func (接收者变量 接收者类型) 方法名(参数列表) (返回参数) {
    函数体
}

//定义struct(person,dog)
type person struct {
	name string
	age  int
}
type dog struct {
	name string
}

//定义构造函数(newPerson,newDog)
func newPerson(name string, age int) *person {
	return &person{
		name: name,
		age:  age,
	}
}
func newDog(name string) *dog {
	return &dog{
		name: name,
	}
}
// 方法是作用于特定类型的函数
// 接受者表示的是调用该方法的具体类型变量，多用类型名的首字母小写表示
func (d dog) wang() {
	fmt.Printf("%s在旺旺旺\n", d.name)
}

// 使用值接受者：传拷贝进入
func (p person) oneNewYear() {
	p.age++
}

// 使用指针接受者：传内存地址进入
func (p *person) oneNewYearP() {
	p.age++
}

func main() {
	p1 := newPerson("lqx", 18)
	//传拷贝变量进入方法：
  
	p1.oneNewYear()
	fmt.Printf("使用传拷贝变量的方式到方法中的结果：%d\n", p1.age) //使用传拷贝变量的方式到方法中的结果：18

	//传指针进入：
	p1.oneNewYearP()
	fmt.Printf("使用传指针的方式到方法中的结果：%d\n", p1.age) //使用传指针的方式到方法中的结果：19

	d1 := newDog("小狗")
	fmt.Println(d1) //&{wangdog}
	d1.wang()
}
```

### 值接收者和指针接收者的区别

- 值接受者
  1. 使用值接收者的方法不能修改结构体变量
  2. 使用指针接收者的方法可以修改结构体的变量,课上过年长一岁的例子.
- 指针接受者
  1. 需要修改接收者中的值
  2. 接收者是拷贝代价比较大的大对象
  3. 保证一致性，如果有某个方法使用了指针接收者，那么其他的方法也应该使用指针接收者。

## 给自定义类型添加方法

```go
//给自定义类型加方法
// 不能给别的包里面的类型添加方法，只能给自己包里的类型添加方法
type myInt int

func (m myInt) hello() {
	fmt.Println("我重写后的int类型")
}
func main() {
	num := myInt(1000)
	num.hello() //我重写后的int类型
}
```

## 嵌套结构体

### 嵌套匿名结构体

### 嵌套结构体字段名冲突

```go
//结构体的嵌套
type address struct {
	province string
	city     string
}
type workPlace struct {
	province string
	city     string
}
type person struct {
	name      string
	age       int
	address   //匿名嵌套结构体
	workPlace workPlace  //普通嵌套结构体
	// city string
}
type company struct {
	name string
	address
}
func main() {
	//结构体初始化
	p1 := person{
		name: "刘祺祥",
		age:  25,
		address: address{
			province: "山西",
			city:     "太原",
		},
		workPlace: workPlace{
			province: "北京",
			city:     "北京",
		},
	}
	fmt.Println(p1)
	// fmt.Println(p1.city)  //如果匿名结构体中没有重名的字段，那么会先在自己结构体找，如果没有，才会去匿名结构体中查找
	fmt.Println(p1.name, p1.address.city, p1.workPlace.city)
}
```

## 结构体的继承

```go

//定义一个动物类struct
type animal struct {
	name string
}
//给一个动物类添加一个移动的方法
func (a animal) move() {
	fmt.Printf("%s还会动\n", a.name)
}

//定义一个dog的struct
type dog struct {
	foot   uint8
	animal //animal拥有的方法，dog都会得到
}
//给dog实现一个汪汪汪的方法
func (d dog) wang() {
	fmt.Printf("%s，它有%d条腿\n", d.name, d.foot)
}
func main() {
	d1 := dog{
		foot: 4,
		animal: animal{
			name: "大黄狗",
		},
	}
	d1.wang()
	d1.move() //d1调用animal的方法，实现继承的关系
}
```

## 结构体字段的可见性

结构体中字段大写开头表示可公开访问，小写表示私有（仅在定义当前结构体的包中可访问）。

## 结构体与json序列化

### 结构体的标签（tag）

```go
// 标签：	Name string `json:"name" db:"name" ini:"name"`

//struct中的值首字母都必须大写
//1.序列化： 把go语言中的结构体变量 --- > json格式的字符串
//2.反序列化： 把json格式的字符串 --- >go语言中能够识别的结构体变量
type person struct {
	Name string `json:"name" db:"name" ini:"name"` //反撇代表的是打标签，在json中这个Name为小写的"name"
	Age  int    `json:"age"`
}

func main() {
	p1 := person{
		Name: "liuqixiang",
		Age:  20,
	}
	//序列化
	b, err := json.Marshal(p1) //返回的是一个[]byte类型
	if err != nil {
		fmt.Printf("marshal failed，err:%v\n", err)
	}
	fmt.Printf("%v\n", string(b)) //byte转换为string  {"name":"liuqixiang","age":20}

	//反序列化
	jsonStr := `{"name":"liuqixiang","age":20}`
	var p2 person
	json.Unmarshal([]byte(jsonStr), &p2) //传指针是为了能在json.Unmarshal内部修改p2的值
	fmt.Println(p2)
}

```

### 序列化注意事项

1. 序列化的时候需要结构体内部的值首字母大写
2. 反序列化需要传指针

## 链表的定义

> 每个节点包含下一个节点的地址，这样把所有的节点串起来，通常把链表中的第一个节点叫做链表头

### 链表尾部插入法

```go
package main

import "fmt"

type Student struct {
	Name  string
	Age   int
	score float32
	next  *Student
}
// 链表尾部插入法
func student3() {
	//设置链表头
	var head Student  //指定一个头部结构体
	head.Name = "hua"
	head.Age = 10
	head.score = 100

	var stu1 Student
	stu1.Name = "stu1_hua"
	stu1.Age = 20
	stu1.score = 22

	var stu2 Student
	stu2.Name = "stu2_hua"
	stu2.Age = 30
	stu2.score = 23
	stu1.next = &stu2 //stu1 结构体的后面是stu2

	var stu3 Student
	stu3.Name = "stu3_hua"
	stu3.Age = 40
	stu3.score = 24
	stu2.next = &stu3

	head.next = &stu1

	var p *Student = &head
	for p != nil {
		fmt.Println(*p)
		p = p.next
		// {hua 10 100 0xc00008e030}
		// {stu1_hua 20 22 0xc00008e000}
		// {stu2_hua 30 23 <nil>}
	}
}
```

