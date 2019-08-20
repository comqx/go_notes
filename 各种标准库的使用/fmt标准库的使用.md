- [print](#print)
	- [Print](#print)
	- [Println](#println)
	- [Printf](#printf)
		- [普通用法（General）](#普通用法general)
		- [布尔类型（Boolean）](#布尔类型boolean)
		- [整数（Integer）](#整数integer)
		- [整数宽度（Integer width）](#整数宽度integer-width)
		- [浮点型（Float）](#浮点型float)
		- [字符串（String）](#字符串string)
		- [字符串宽度（String Width）](#字符串宽度string-width)
		- [结构（Struct）](#结构struct)
		- [指针（Pointer）](#指针pointer)
		- [例子](#例子)
- [sprint](#sprint)
	- [Sprintf](#sprintf)
- [fprint](#fprint)
	- [Fprintln](#fprintln)
	- [Fprintf](#fprintf)
- [Scan](#scan)
	- [Scan](#scan-1)
	- [Scanf](#scanf)
	- [Scanln](#scanln)
# print
## Print

>就是一般的标准输出，但是不换行

## Println
>可以打印出字符串，和变量，自动添加换行

## Printf
>只可以打印出格式化的字符串，可以输出字符串类型的变量，不可以输出整形变量和整形

也就是说，当需要格式化输出信息时一般使用Printf,其他时候都使用Println
### 普通用法（General）
```
%v 以默认的方式打印变量的值
%#v 相应值的go语法表示
%T 打印变量的类型
%% 字面上的%,不是值的占位符
```
### 布尔类型（Boolean）
```
%t 打印true或false
```
### 整数（Integer）
```go
%+d 带符号的整型，fmt.Print("%+d",255),输出是+255
%q 打印单引号

%b 打印整型的二进制
%o 不带零的八进制
%#o 带零的八进制
%x 小写的十六进制
%X 大写的十六进制
%#x 带0x的十六进制

%U 打印Unicode字符
%#U 打印带字符的Unicode
	var exint int = 23567
	fmt.Printf("unicode:%#U\n", exint) //将int类型转换为unicode字符，并且携带对应的字符 //unicode:U+5C0F '小'

```
### 整数宽度（Integer width）
```
%5d 表示该整型最大长度是5，下面这段代码
	var exint int = 100
	fmt.Printf("5d:%5d\n", exint) //5d:  100  //在整数前面补充空格，总长度为5
	
%-5d则相反，打印结果会自动左对齐
%05d会在数字前面补零。
	fmt.Printf("5d:%05d\n", exint) //5d:00100
```
### 浮点型（Float）

```
%f (=%.6f) 6位小数点
%e (=%.6e) 6位小数点（科学计数法）
%g 用最少的数字来表示
%.3g 最多3位数字来表示
%.3f 最多3位小数来表示
	var exflo float64 = 3000.14152341231324567876
	fmt.Printf("float:%.3f\n", exflo) //保留3位小数  float:3000.142
	fmt.Printf("float:%.3g\n", exflo) //带整数位，保留3个数字 float:3e+03
	
```
### 字符串（String）
```go
%s 正常输出字符串
%q 字符串带双引号，字符串中的引号带转义符
%#q 字符串带反引号，如果字符串内有反引号，就用双引号代替
%x 将字符串转换为小写的16进制格式
%X 将字符串转换为大写的16进制格式
% x 带空格的16进制格式
%c acsii码转换为字符

	fmt.Printf("string:%X\n", exstr)  //将字符串转换为大写的16进制格式 //string:6C697571697869616E67EFBC8CE58898E7A5BAE7A5A5
	fmt.Printf("string:%x\n", exstr)  // 将字符串转换为小写的16进制格式  //string:6c697571697869616e67efbc8ce58898e7a5bae7a5a5
	fmt.Printf("string:% x\n", exstr) //将字符串转换为小写的带有空格的16进制格式 //string:6c 69 75 71 69 78 69 61 6e 67 ef bc 8c e5 88 98 e7 a5 ba e7 a5 a5
	fmt.Printf("acsii:%c\n", exint)   //将acsii码转化为对于的值 //acsii:d
```
### 字符串宽度（String Width）
```
%5s 最小宽度为5
%-5s 最小宽度为5（左对齐）
%.5s 最大宽度为5
%5.7s 最小宽度为5，最大宽度为7
%-5.7s 最小宽度为5，最大宽度为7（左对齐）
%5.3s 如果宽度大于3，则截断
%05s 如果宽度小于5，就会在字符串前面补零
```
### 结构（Struct）
```
%v 正常打印。比如：{sam {12345 67890}}
%+v 带字段名称。比如：{name:sam phone:{mobile:12345 office:67890}
%#v 用Go的语法打印。
比如main.People{name:”sam”, phone:main.Phone{mobile:”12345”, office:”67890”}}
```

### 指针（Pointer）
```
%p 带0x的指针
%#p 不带0x的指针
```
### 例子
```go
package main

import "fmt"

func main(){
	var a int
	var b bool
	c := 'a'

	fmt.Printf("%+v\n",a) //默认打印，可以打印结构体，有+号会添加字段名
	fmt.Printf("%#v\n",b) //go语法表示
	fmt.Printf("%T\n",c) //打印值类型
	fmt.Printf("90%%\n") //打印%
	fmt.Printf("%t\n",b) //单词true或者false
	fmt.Printf("%b\n",100) //打印无小数部分
	fmt.Printf("%f\n",100.221) //打印小数
	fmt.Printf("%q\n","this is test") //双引号围绕的字符串,go语法安全转义
	fmt.Printf("%x\n",42145342) //十六进制，小写字母，每字节俩个字符
	fmt.Printf("%p\n",&a)  //打印内存地址

	str := fmt.Sprintf("a=%d",a)   //int转换str类型
	fmt.Printf("%q\n",str)
}
```
# sprint
## Sprintf

>是输出到串，一般是直接申请输出到一个字符串中，可以用来将大量数字数据转换成字符串
```go
	var a = "hello"
	var b = "wolrd"
	c := fmt.Sprintf("%s+%s", a, b)
	fmt.Println(c) //hello+wolrd
```
# fprint
## Fprintln

```go
// 向标准输出写入内容
fmt.Fprintln(输入句柄, 内容)

fmt.Fprintln(os.Stdout, "向标准输出写入内容")
```

## Fprintf

>是输出到文件，当然，这个文件也可能是虚拟的文件

```GO
	fileobj, err := os.OpenFile("./xx.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) //定义了文件名，文件权限
	if err != nil {
		fmt.Println("打开文件出错，err:", err)
		return
	}

	name := "沙河小王子"
	fmt.Fprintf(fileobj, "往文件中写入信息：%s", name) //写入内容到文件中，
```

# Scan

## Scan

> scan从标准输入扫描文本,读取由空白符分割的值保存到传递给对应的地址

```go
var (
	name string
  age int
  married bool
)
fmt.Scan(&name,&age,&married)
fmt.Printf("扫描结果 name:%s age:%d married:%t \n", name, age, married)
```

## Scanf

```go
	var (
		name    string
		age     int
		married bool
	)
	fmt.Scanf("1:%s 2:%d 3:%t", &name, &age, &married) //1:abc 2:10 3:false 格式化接收值并赋值给指定的变量
	fmt.Printf("扫描结果 name:%s age:%d married:%t \n", name, age, married)
```

## Scanln

# Errorf

`Errorf`函数根据format参数生成格式化字符串并返回一个包含该字符串的错误。

```go
func Errorf(format string, a ...interface{}) error

//自定义错误输出
err := fmt.Errorf("这是一个错误")
```

# bufio.NewReader

有时候我们想完整获取输入的内容，而输入的内容可能包含空格，这种情况下可以使用`bufio`包来实现.

```go
func bufioDemo() {
	reader := bufio.NewReader(os.Stdin) // 从标准输入生成读对象
	fmt.Print("请输入内容：")
	text, _ := reader.ReadString('\n') // 读到换行
	text = strings.TrimSpace(text)
	fmt.Printf("%#v\n", text)
}
```