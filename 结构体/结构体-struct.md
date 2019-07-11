[TOC]

# struct

- 用来自定义复杂数据结构
- struct里面可以包含多个字段(属性)
- struct类型可以定义方法，注意和函数的区分
- struct类型是值类型
- struct类型可以嵌套
- go语言没有class类型，只有struct类型

## struct的声明

```go
type 标识符 struct{
  field1 type
  filed2 type
}
```

## struct的定义

> Struct 中字段访问，和其他语言一样，使用点

```go
var stu student
stu.name = "tony"
stu.age = 18
stu.score = 20
fmt.Printf("name=%s age=%s score=%s",stu.name,stu.age,stu.score)
```



### struct定义的3种形式

```go
a. var stu Student
b. var stu *Student = new(Student) //new返回一个指针
c. var stu *Student = &Student{}  //定义一个结构体的指针
```

1. 其中b,c返回的都是指向结构体的指针，访问形式如下：

   `Stu.name 、stu.age和stu.score或者(*stu).name 、(*stu).age等`

## struct的初始化

> struct的内存布局， struct中的所有字段的内存是连续的，布局如下

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

