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