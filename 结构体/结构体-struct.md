[TOC]

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
type person struct{
  name，city string
  age int8
}
结构体用来描述一组值，比如一个人有名字，年龄，居住城市等
```

## struct的实例化

> - struct的内存布局， struct中的所有字段的内存是连续的
>
> - 结构体只有实例化以后，才会真正分配内存，也就是必须实例化以后才能使用结构体的字段
>
> - 结构体本身也是一种类型，可以像声明内置类型一样使用var关键字声明结构体类型

```go
var 结构体实例 结构体类型

type person struct {
	name string
	city string
	age  int8
}
func main() {
	var p1 person
	p1.name = "沙河娜扎"
	p1.city = "北京"
	p1.age = 18
	fmt.Printf("p1=%v\n", p1)  //p1={沙河娜扎 北京 18}
	fmt.Printf("p1=%#v\n", p1) //p1=main.person{name:"沙河娜扎", city:"北京", age:18}
}
```

## 匿名结构体

```go
    var user struct{Name string; Age int} //定义匿名结构体
    user.Name = "小王子"
    user.Age = 18
    fmt.Printf("%#v\n", user)
```

### 链表的定义

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

